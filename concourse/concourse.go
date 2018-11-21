package concourse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Concourse struct {
	pathToFly string
	target    string
}

func New(pathToFly, target string) Concourse {
	return Concourse{pathToFly: pathToFly, target: target}
}

func (c Concourse) BuildOutput(bid int) (string, error) {
	cmd := exec.Command(c.pathToFly, "--target", c.target, "watch", "--build", strconv.Itoa(bid))

	var stdOutErr bytes.Buffer
	cmd.Stdout = &stdOutErr
	cmd.Stderr = &stdOutErr
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err.Error(), stdOutErr.String()))
	}

	return stdOutErr.String(), nil
}

func (c Concourse) Builds(pipelineName string) ([]Build, error) {
	cmd := exec.Command(c.pathToFly, "--target", c.target, "builds", "--pipeline", pipelineName, "--json")

	var stdOutErr bytes.Buffer
	cmd.Stdout = &stdOutErr
	cmd.Stderr = &stdOutErr
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err.Error(), stdOutErr.String()))
	}

	var builds []Build
	err = json.Unmarshal(stdOutErr.Bytes(), &builds)
	if err != nil {
		panic(err)
	}

	return builds, nil
}

func (c Concourse) Pipelines() ([]Pipeline, error) {
	cmd := exec.Command(c.pathToFly, "--target", c.target, "pipelines", "--all", "--json")

	var stdOutErr bytes.Buffer
	cmd.Stdout = &stdOutErr
	cmd.Stderr = &stdOutErr
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err.Error(), stdOutErr.String()))
	}

	var pipelines []Pipeline
	err = json.Unmarshal(stdOutErr.Bytes(), &pipelines)
	if err != nil {
		panic(err)
	}

	return pipelines, nil
}

func (c Concourse) TargetInfo() (TargetInfo, error) {
	cmd := exec.Command(c.pathToFly, "targets")

	var stdOutErr bytes.Buffer
	cmd.Stdout = &stdOutErr
	cmd.Stderr = &stdOutErr
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err.Error(), stdOutErr.String()))
	}

	lines := strings.Split(stdOutErr.String(), "\n")

	var targetInfo TargetInfo
	for _, l := range lines {
		allParts := strings.Fields(l)
		dateTime := strings.Join(allParts[3:], " ")
		if allParts[0] == c.target {
			expiryTime, err := time.Parse(time.RFC1123, dateTime)
			if err != nil {
				panic(err)
			}

			targetInfo = TargetInfo{Name: allParts[0], URL: allParts[1], Team: allParts[2], Expiry: expiryTime}
			return targetInfo, nil
		}
	}

	return TargetInfo{}, errors.New("did not find target")
}
