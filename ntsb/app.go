package ntsb

import (
	"fmt"
	"io"
	"path"
	"regexp"

	"github.com/concourse/atc"
	"github.com/concourse/atc/event"
	"github.com/concourse/fly/rc"
	"github.com/concourse/go-concourse/concourse"
)

type App struct {
	flagSet      flagSet
	rcLoadTarget func(target rc.TargetName, tracing bool) (rc.Target, error)
}

func NewApp(flagSet flagSet, rcLoadTarget func(target rc.TargetName, tracing bool) (rc.Target, error)) App {
	return App{
		flagSet:      flagSet,
		rcLoadTarget: rcLoadTarget,
	}
}

func (a App) Run(args []string) error {
	var (
		jobName    string
		pattern    string
		buildCount int
	)

	a.flagSet.StringVar(&jobName, "j", "", "name of concourse job")
	a.flagSet.StringVar(&pattern, "p", "", "pattern to search in build log")
	a.flagSet.IntVar(&buildCount, "c", 200, "how many builds to search")
	a.flagSet.Parse(args)

	if jobName == "" {
		return fmt.Errorf("job name is required\n")
	}

	if pattern == "" {
		return fmt.Errorf("pattern is required\n")
	}

	if buildCount < 0 {
		return fmt.Errorf("count must be a positive integer\n")
	}

	const verbose = false
	target, err := a.rcLoadTarget("releng", verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	page := concourse.Page{Limit: buildCount}

	team := target.Team()
	pipelines, err := team.ListPipelines()
	if err != nil {
		return fmt.Errorf("pipelines: %s", err)
	}

	// get all failing builds
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

	// collect builds with logs containing specified pattern
	client := target.Client()
	buildsWithPattern := []atc.Build{}

	for _, b := range allFailingBuilds {
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
				matched, err := regexp.MatchString(pattern, e.Payload)
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

	for _, b := range buildsWithPattern {
		fmt.Println(path.Join(target.URL(), "teams", b.TeamName, "pipelines", b.PipelineName,
			"jobs", b.JobName, "builds", b.Name))
	}

	return nil
}
