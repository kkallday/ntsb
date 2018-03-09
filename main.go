package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/concourse/atc"
	"github.com/concourse/fly/rc"
	"github.com/concourse/go-concourse/concourse"
)

func main() {
	var jobName string
	flag.StringVar(&jobName, "j", "", "name of concourse job")
	flag.Parse()

	err := mainWithError(jobName)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

func mainWithError(requestedJobName string) error {
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
			requestedJobName,
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

	for _, b := range allFailingBuilds {
		log.Printf("failure! %+v\n", b)
	}

	return nil
}
