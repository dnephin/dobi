package cmd

import (
	"fmt"
	"github.com/dnephin/dobi/config"
	"github.com/spf13/cobra"
)

func newListCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}
	return cmd
}

func runList(opts *dobiOptions) error {
	conf, err := config.Load(opts.filename)
	if err != nil {
		return err
	}

	printTasks(conf)
	return nil
}

func printTasks(config *config.Config) {
	for _, name := range config.Sorted() {
		fmt.Printf("  %-20s %s\n", name, config.Resources[name])
	}
}
