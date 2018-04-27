package main

import (
	"log"
	"os"

	"github.com/kkallday/ntsb/ntsb"
)

func main() {
	app := ntsb.NewApp()
	err := app.Run(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
