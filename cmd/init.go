package cmd

import (
	"github.com/spf13/cobra"
)


func newInitCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run the remove action for all resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(
		newGolangCommand(opts),
	)
	return cmd
}
