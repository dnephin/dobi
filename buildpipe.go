package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
	"github.com/dnephin/buildpipe/steps"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filename = kingpin.Flag("filename", "Path to config file").String()
	verbose  = kingpin.Flag("verbose", "Verbose logging").Bool()
)

func initLogging(level log.Level) {
	log.SetOutput(os.Stderr)
	log.SetLevel(level)
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// TODO: verbose flag
	initLogging(log.InfoLevel)

	conf, err := config.Load(*filename)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}
	if err := steps.Run(conf); err != nil {
		log.Fatalf("%s", err)
	}
}
