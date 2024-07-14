package cmd

import (
	"github.com/Andreagit97/bpftree/pkg/task"
	"github.com/spf13/cobra"
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
