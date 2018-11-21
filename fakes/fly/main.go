package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if os.Args[1] == "targets" {
		fmt.Printf(`name      url                  team       expiry                       
target-1  https://example.com  some-team  Thu, 25 Oct 2018 06:22:59 UTC
target-2  https://example.org  main       Fri, 16 Nov 2018 18:04:37 UTC
`)
		return
	}

	switch os.Args[3] {
	case "pipelines":
		fmt.Println(pipelinesJSON)
	case "builds":
		fmt.Println(getBuilds(os.Args[5]))
	case "watch":
		bid, err := strconv.Atoi(os.Args[5])
		if err != nil {
			panic(err)
		}
		fmt.Println(getBuildOutput(bid))
	default:
		panic("Unexpected command: " + strings.Join(os.Args, " "))
	}
}

func getBuilds(pipelineName string) string {
	switch pipelineName {
	case "cool-pipeline":
		return coolBuildsJSON
	case "boring-pipeline":
		return boringBuildsJSON
	default:
		panic("Unexpected pipeline name: " + pipelineName)
	}
}

func getBuildOutput(bid int) string {
	switch bid {
	case 9817:
		fallthrough
	case 3781:
		fallthrough
	case 47187:
		return "blahajkaljlakthis-is-a-patternkajsdlkfj"
	case 7788:
		return "fjkjfalkjdlknothingmatchesjdakjfkl"
	default:
		panic(fmt.Sprintf("Unexpected build id: %d", bid))
	}
}
