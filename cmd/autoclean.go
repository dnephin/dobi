package cmd

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks"
	"github.com/spf13/cobra"
)

func newCleanCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoclean",
		Short: "Run the remove action for all resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClean(opts)
		},
	}
	return cmd
}

func runClean(opts *dobiOptions) error {
	conf, err := config.Load(opts.filename)
	if err != nil {
		return err
	}

	client, err := buildClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %s", err)
	}

	return tasks.Run(tasks.RunOptions{
		Client: client,
		Config: conf,
		Tasks:  removeTasks(conf),
		Quiet:  opts.quiet,
	})
}

func removeTasks(conf *config.Config) []string {
	resources := conf.Sorted()
	tasks := []string{}
	for i := len(resources) - 1; i >= 0; i-- {
		tasks = append(tasks, resources[i]+":rm")
	}
	return tasks
}
