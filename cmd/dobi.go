package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks"
	"github.com/dnephin/dobi/tasks/client"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

const (
	// DefaultDockerAPIVersion is the default version of the docker API to use
	DefaultDockerAPIVersion = "1.23"
)

var (
	version   = "0.7"
	gitsha    = "unknown"
	buildDate = ""
)

type dobiOptions struct {
	filename string
	verbose  bool
	quiet    bool
	tasks    []string
	version  bool
}

// NewRootCommand returns a new root command
func NewRootCommand() *cobra.Command {
	var opts dobiOptions

	cmd := &cobra.Command{
		Use:                   "dobi [flags] RESOURCE[:ACTION] [RESOURCE[:ACTION]...]",
		Short:                 "A build automation tool for Docker applications",
		SilenceUsage:          true,
		SilenceErrors:         true,
		TraverseChildCommands: true,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.tasks = args
			return runDobi(opts)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initLogging(opts.verbose, opts.quiet)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.filename, "filename", "f", "dobi.yaml", "Path to config file")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Quiet")
	flags.BoolVar(&opts.version, "version", false, "Print version and exit")

	flags.SetInterspersed(false)
	cmd.AddCommand(newListCommand(&opts))
	return cmd
}

func runDobi(opts dobiOptions) error {
	if opts.version {
		printVersion()
		return nil
	}

	conf, err := config.Load(opts.filename)
	if err != nil {
		return err
	}

	client, err := buildClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %s", err)
	}

	return tasks.Run(tasks.RunOptions{
		Client: client,
		Config: conf,
		Tasks:  opts.tasks,
		Quiet:  opts.quiet,
	})
}

func initLogging(verbose, quiet bool) {
	logger := logging.Log
	if verbose {
		logger.Level = log.DebugLevel
	}
	if quiet {
		logger.Level = log.WarnLevel
	}
	logger.Out = os.Stderr

	formatter := &logging.Formatter{}
	log.SetFormatter(formatter)
	logger.Formatter = formatter
}

func buildClient() (client.DockerClient, error) {
	apiVersion := os.Getenv("DOCKER_API_VERSION")
	if apiVersion == "" {
		apiVersion = DefaultDockerAPIVersion
	}
	// TODO: args for client
	client, err := docker.NewVersionedClientFromEnv(apiVersion)
	if err != nil {
		return nil, err
	}
	log.Debug("Docker client created")
	return client, nil
}

func printVersion() {
	fmt.Printf("dobi version %v (build: %v, date: %s)\n", version, gitsha, buildDate)
}
