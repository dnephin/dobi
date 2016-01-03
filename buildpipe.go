package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
	"github.com/dnephin/buildpipe/tasks"
	flag "github.com/spf13/pflag"
)

var (
	filename = flag.StringP("filename", "f", "buildpipe.yaml", "Path to config file")
	verbose  = flag.BoolP("verbose", "v", false, "Verbose")
	quiet    = flag.BoolP("quiet", "q", false, "Quiet")
)

func initLogging(verbose, quiet bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	if quiet {
		log.SetLevel(log.WarnLevel)
	}
	log.SetOutput(os.Stderr)
}

func main() {
	cmd := flag.CommandLine
	cmd.Init(os.Args[0], flag.ExitOnError)
	cmd.SetInterspersed(false)
	flag.Usage = func() {
		out := os.Stderr
		fmt.Fprintf(out, "Usage:\n  %s [OPTIONS] PIPELINE... \n\n", os.Args[0])
		fmt.Fprintf(out, "Options:\n")
		cmd.PrintDefaults()
	}
	flag.Parse()
	initLogging(*verbose, *quiet)

	pipelines := flag.Args()

	conf, err := config.Load(*filename)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}
	if err := tasks.Run(conf, pipelines); err != nil {
		log.Fatalf("%s", err)
	}
}
