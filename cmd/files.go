package cmd

import (
	"github.com/Andreagit97/bpftree/pkg/task"
	"github.com/spf13/cobra"
)

var (
	// todo! add description
	filesCmd = &cobra.Command{
		Use:     "files <field-name> <field-value>",
		Aliases: []string{"fds"},
		Short:   "Shows all files of a certain thread",
		Long:    "Shows all files of a certain thread",
		Example: " - bpftree files tid 1 -> print all files for task with tid=1",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			task.TasksFiles(args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(filesCmd)
}
