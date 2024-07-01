package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	fieldsCmd = &cobra.Command{
		Use:     "fields",
		Aliases: []string{"f"},
		Short:   "Shows all available fields",
		Run: func(cmd *cobra.Command, args []string) {
			// todo!: this shouldn't stay in the render package
			fmt.Println("not supported")
			os.Exit(1)

			// task.RenderAllowedFields()
		},
	}
)

func init() {
	rootCmd.AddCommand(fieldsCmd)
}
