package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if os.Args[1] == "targets" {
		fmt.Printf(`target-1    https://example.com    some-team    Thu, 25 Oct 2018 06:22:59 UTC
target-2    https://example.org    main         Fri, 16 Nov 2018 18:04:37 UTC
`)
		return
	}

	switch os.Args[3] {
	case "builds":
		fmt.Printf(`[
	{
		"id": 1234,
		"name": "2",
		"status": "succeeded",
		"job_name": "job-a"
	},
	{
		"id": 4321,
		"name": "3",
		"status": "failed",
		"job_name": "job-b"
	}
]
`)
	case "pipelines":
		fmt.Printf(`[
	{
		"id": 450,
		"name": "pipeline-a",
		"team_name": "some-team-1"
	},
	{
		"id": 451,
		"name": "pipeline-b",
		"team_name": "some-team-1"
	},
	{
		"id": 943,
		"name": "pipeline-c",
		"team_name": "some-team-2"
	}
]
`)
	case "watch":
		fmt.Printf("something\nsomething else\neven more\n")
	default:
		panic("Unexpected command: " + strings.Join(os.Args, " "))
	}
}
