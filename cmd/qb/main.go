package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/winebarrel/qb"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
}

func main() {
	flags := parseFlags()
	task, err := qb.NewTask(&flags.TaskOpts, &flags.RecorderOpts)

	if err != nil {
		log.Fatalf("failed to build task: %s", err)
	}

	if flags.initialize {
		err = task.Initialize()

		if err != nil {
			log.Fatalf("failed to initialize: %s", err)
		}

		return
	}

	err = task.Prepare()

	if err != nil {
		log.Fatalf("failed to prepare task: %s", err)
	}

	rec, err := task.Run()

	if err != nil {
		log.Fatalf("failed to run task: %s", err)
	}

	err = task.Close()

	if err != nil {
		log.Fatalf("failed to close task: %s", err)
	}

	if !flags.OnlyPrint {
		report := rec.Report()
		rawJson, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(rawJson))
	}
}
