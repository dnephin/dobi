package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

const (
	// DockerAPIVersion is the version of the docker API to use
	DockerAPIVersion = "1.23"
)

type dobiOptions struct {
	filename string
	list     bool
	verbose  bool
	quiet    bool
	tasks    []string
}

// NewRootCommand returns a new root command
func NewRootCommand() *cobra.Command {
	var opts dobiOptions

	cmd := &cobra.Command{
		Use:           "dobi TASK [TASK...]",
		Short:         "A build tool for Docker application.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.tasks = args
			return runDobi(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.filename, "filename", "f", "dobi.yaml", "Path to config file")
	flags.BoolVar(&opts.list, "list", false, "List all available tasks")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Quiet")

	flags.SetInterspersed(false)
	return cmd
}

func runDobi(opts dobiOptions) error {
	initLogging(opts.verbose, opts.quiet)

	conf, err := config.Load(opts.filename)
	if err != nil {
		return fmt.Errorf("Failed to load config from %q: %s", opts.filename, err)
	}

	// TODO: make this a Command instead of a flag
	if opts.list {
		listTasks(conf)
		return nil
	}

	client, err := buildClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %s", err)
	}

	return tasks.Run(tasks.RunOptions{Client: client, Config: conf, Tasks: opts.tasks})
}

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
	client, err := docker.NewVersionedClientFromEnv(DockerAPIVersion)
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
