package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Andreagit97/bpftree/pkg/task"
)

var (
	fieldsCmd = &cobra.Command{
		Use:     "fields",
		Aliases: []string{"f"},
		Short:   "Shows all available fields",
		Run: func(cmd *cobra.Command, args []string) {
			task.RenderAllowedFields()
		},
	}
)

func init() {
	rootCmd.AddCommand(fieldsCmd)
}
