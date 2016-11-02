package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"log"
)

var extensions []string

func newInitCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run the remove action for all resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := EnforceArgumentExistance()
			if err != nil {
				return err
			}

			switch os.Args[2] {
			case "golang":
				err := setGolangDefault()
				return err
			default:
				return fmt.Errorf("%s language is not supported by dobi yet", os.Args[2])
			}

			log.Println("here at the end")

			return nil
		},
	}
	return cmd
}

func EnforceArgumentExistance() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("need another argument")
	}
	return nil
}

func setGolangDefault() error {
	// CHECK IF GO EXIST
	os.lo

	// CHECK IF VENDOR FILE EXIST

	return nil
}