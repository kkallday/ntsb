package ntsb

import (
	"fmt"
	"regexp"

	concoursefly "github.com/kkallday/ntsb/concourse"
)

type App struct {
	concourse concourse
	logger    logger
}

//go:generate counterfeiter . concourse
type concourse interface {
	BuildOutput(bid int) (string, error)
	Builds(pipelineName string) ([]concoursefly.Build, error)
	TargetInfo() (concoursefly.TargetInfo, error)
	Pipelines() ([]concoursefly.Pipeline, error)
}

//go:generate counterfeiter . logger
type logger interface {
	Println(...interface{})
}

func NewApp(concourse concourse, logger logger) App {
	return App{
		concourse: concourse,
		logger:    logger,
	}
}

func (a App) Run(pattern string) error {
	targetInfo, err := a.concourse.TargetInfo()
	if err != nil {
		panic(err)
	}

	/*
		if targetInfo.Expired {
			panic(errors.New("session expired, please login"))
		}
	*/

	pipelines, err := a.concourse.Pipelines()
	if err != nil {
		panic(err)
	}

	var urlsToMatchingBuilds []string

	for _, pipeline := range pipelines {
		builds, err := a.concourse.Builds(pipeline.Name)
		if err != nil {
			panic(err)
		}

		for _, b := range builds {
			if b.Status != "failed" && b.Status != "errored" {
				continue
			}

			buildOutput, err := a.concourse.BuildOutput(b.ID)
			if err != nil {
				panic(err)
			}

			matched, err := regexp.MatchString(pattern, buildOutput)
			if err != nil {
				panic(err)
			}

			if !matched {
				continue
			}

			url := fmt.Sprintf("%s/teams/%s/pipelines/%s/jobs/%s/builds/%s", targetInfo.URL, b.TeamName, pipeline.Name, b.JobName, b.Name)
			urlsToMatchingBuilds = append(urlsToMatchingBuilds, url)
		}
	}

	a.logger.Println("The following build(s) contain text matching the pattern:")

	for _, url := range urlsToMatchingBuilds {
		a.logger.Println(url)
	}

	return nil
}
