package main

import (
	"log"

	"github.com/dnephin/dbp/config"
	"github.com/dnephin/dbp/steps"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filename = kingpin.Flag("filename", "Path to config file").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	conf, err := config.Load(*filename)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}
	if err := steps.Run(conf); err != nil {
		log.Fatalf("Failed to run steps:\n%s", err)
	}
}
