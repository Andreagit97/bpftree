package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Andreagit97/bpftree/pkg/task"
)

var (
	filesCmd = &cobra.Command{
		Use:     "files <thread-id>",
		Aliases: []string{"fds"},
		Short:   "Shows all files of a certain thread",
		Long:    "Shows all files of a certain thread",
		Example: " - bpftree files 1 -> print all files for task with tid=1",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			task.PrintTaskFiles(args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(filesCmd)
}
