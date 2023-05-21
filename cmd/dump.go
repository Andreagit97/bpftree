package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Andreagit97/bpftree/pkg/task"
)

var (
	dumpCmd = &cobra.Command{
		Use:     "dump <output-file-name>",
		Aliases: []string{"d"},
		Short:   "Dump system tasks into a file",
		Long: `Dump system tasks into a file.
The first argument is the file used for the dump.
This file can then be read by bpftree using the '--capture' flag`,
		Example: `  - bpftree dump outfile.tree -> dump tasks into outfile.tree`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := task.DumpTasksIntoFile(args[0]); err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}
