package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	concoursefly "github.com/kkallday/ntsb/concourse"
	"github.com/kkallday/ntsb/ntsb"
)

func main() {
	var (
		target  string
		pattern string
	)

	fs := flag.NewFlagSet("ntsb", flag.ExitOnError)
	fs.StringVar(&target, "target", "", "fly target to use for authentication")
	fs.StringVar(&pattern, "pattern", "", "pattern to search (regex)")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf(err.Error())
	}

	if target == "" {
		panic("target is a required argument")
	}

	if pattern == "" {
		panic("pattern is a required argument")
	}

	pathToFly, err := exec.LookPath("fly")
	if err != nil {
		panic(err)
	}
	concourse := concoursefly.New(pathToFly, target)
	logger := ntsb.NewLogger()
	app := ntsb.NewApp(concourse, logger)

	err = app.Run(pattern)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
