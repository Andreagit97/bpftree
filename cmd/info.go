package cmd

import (
	"github.com/Andreagit97/bpftree/pkg/task"
	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use:     "info <field-name> <field-value>",
		Aliases: []string{"i"},
		Short:   `Show all task info that match a certain field value`,
		Long: `Show all task info that match a certain field value.
The field name is provided as first argument while the value is provided as second.
The default format used to print a task info is:
[comm] tid: ..., pid: ..., rptid: ..., rppid: ...
You can customize this format using the '--format' flag`,
		Example: `  - bpftree info tid 1 -> print info for tasks with tid=1
  - bpftree i t 1 -> print info for tasks with tid=1 (short form)
  - bpftree info comm systemd -> print info for tasks with comm=systemd
  - bpftree info comm systemd -f 't,p,r' -> print formatted tasks with comm=systemd`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			task.TasksInfo(args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(infoCmd)
}
