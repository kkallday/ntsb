package main

import (
	"flag"
	"log"
	"os"

	"github.com/kkallday/ntsb/ntsb"
)

func main() {
	flagSet := flag.NewFlagSet("ntsb", flag.ExitOnError)

	app := ntsb.NewApp(flagSet)
	err := app.Run(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
