package cmd

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/spf13/cobra"
)

type listOptions struct {
	all bool
}

func newListCommand(opts *dobiOptions) *cobra.Command {
	var listOpts listOptions
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts, listOpts)
		},
	}
	cmd.Flags().BoolVarP(
		&listOpts.all, "all", "a", false,
		"List all resources, including those without descriptions")
	return cmd
}

func runList(opts *dobiOptions, listOpts listOptions) error {
	conf, err := config.Load(opts.filename)
	if err != nil {
		return err
	}

	lines := getDescriptions(conf, listOpts)
	if len(lines) == 0 {
		logging.Log.Warn("No resource descriptions")
		return nil
	}
	fmt.Print(strings.Join(lines, ""))
	return nil
}

func getDescriptions(config *config.Config, listOpts listOptions) []string {
	lines := []string{}
	for _, name := range config.Sorted() {
		res := config.Resources[name]
		desc := res.Describe()
		if desc == "" {
			if !listOpts.all {
				continue
			}
			desc = res.String()
		}
		lines = append(lines, fmt.Sprintf("  %-20s %s\n", name, desc))
	}
	return lines
}
