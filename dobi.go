package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks"
	"github.com/fsouza/go-dockerclient"
	flag "github.com/spf13/pflag"
)

var (
	filename = flag.StringP("filename", "f", "dobi.yaml", "Path to config file")
	list     = flag.Bool("list", false, "List all available tasks")
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

func buildClient() (*docker.Client, error) {
	// TODO: args for client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	log.Info("Docker client created")
	return client, nil
}

func listTasks(config *config.Config) {
	for _, name := range config.Sorted() {
		fmt.Printf("  %-20s %s\n", name, config.Resources[name])
	}
}

func main() {
	cmd := flag.CommandLine
	cmd.Init(os.Args[0], flag.ExitOnError)
	cmd.SetInterspersed(false)
	flag.Usage = func() {
		out := os.Stderr
		fmt.Fprintf(out, "Usage:\n  %s [OPTIONS] TASKS... \n\n", os.Args[0])
		fmt.Fprintf(out, "Options:\n")
		cmd.PrintDefaults()
	}
	flag.Parse()
	initLogging(*verbose, *quiet)

	conf, err := config.Load(*filename)
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	if *list {
		listTasks(conf)
		return
	}

	client, err := buildClient()
	if err != nil {
		log.Fatalf("Failed to create client: %s", err)
	}

	if err := tasks.Run(tasks.RunOptions{
		Client: client,
		Config: conf,
		Tasks:  flag.Args(),
	}); err != nil {
		log.Fatalf(err.Error())
	}
}
