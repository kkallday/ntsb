package main

import (
	"flag"
	"log"
	"os"

	"github.com/concourse/fly/rc"
	"github.com/kkallday/ntsb/ntsb"
)

func main() {
	flagSet := flag.NewFlagSet("ntsb", flag.ExitOnError)
	rcLoadTarget := rc.LoadTarget

	app := ntsb.NewApp(flagSet, rcLoadTarget)
	err := app.Run(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
