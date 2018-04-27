package ntsb

import (
	"fmt"
	"io"
	"regexp"

	"github.com/concourse/atc"
	"github.com/concourse/atc/event"
	"github.com/concourse/fly/rc"
	"github.com/concourse/go-concourse/concourse"
)

type App struct {
	flagSet flagSet
}

type flagSet interface {
	StringVar(p *string, name string, value string, usage string)
	Parse(args []string) error
}

func NewApp(flagSet flagSet) App {
	return App{
		flagSet: flagSet,
	}
}

func (a App) Run(args []string) error {
	var (
		jobName string
		pattern string
	)

	a.flagSet.StringVar(&jobName, "j", "", "name of concourse job")
	a.flagSet.StringVar(&pattern, "p", "", "pattern to search in build log")
	a.flagSet.Parse(args)

	if jobName == "" {
		return fmt.Errorf("job name is required\n")
	}

	if pattern == "" {
		return fmt.Errorf("pattern is required\n")
	}

	const verbose = false
	target, err := rc.LoadTarget("releng", verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	const pageSize = 200
	page := concourse.Page{Limit: pageSize}

	team := target.Team()

	pipelines, err := team.ListPipelines()
	if err != nil {
		return fmt.Errorf("pipelines: %s", err)
	}

	var allFailingBuilds []atc.Build
	for _, p := range pipelines {
		pipelineBuilds, _, found, err := team.JobBuilds(
			p.Name,
			jobName,
			page,
		)
		if err != nil {
			return fmt.Errorf("job-builds: %s", err)
		}

		if !found {
			continue
		}

		for _, build := range pipelineBuilds {
			if build.Status == "failed" {
				allFailingBuilds = append(allFailingBuilds, build)
			}
		}
	}

	client := target.Client()
	buildsWithPattern := []atc.Build{}

	for _, b := range allFailingBuilds {
		matched := false

		eventSource, err := client.BuildEvents(fmt.Sprintf("%d", b.ID))
		if err != nil {
			return err
		}

		for {
			ev, err := eventSource.NextEvent()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return fmt.Errorf("failed to parse next event: %s\n", err)
				}
			}

			if e, ok := ev.(event.Log); ok {
				matched, err = regexp.MatchString(pattern, e.Payload)
				if err != nil {
					return fmt.Errorf("failed to perform regexp: %s\n", err)
				}
				if matched {
					buildsWithPattern = append(buildsWithPattern, b)
					break
				}
			}
		}

		eventSource.Close()
		break
	}

	fmt.Printf("%d builds matched\n", len(buildsWithPattern))

	return nil
}
