package cmd

import (
	"github.com/Andreagit97/bpftree/pkg/task"
	"github.com/spf13/cobra"
)

var (
	lineageCmd = &cobra.Command{
		Use:     "lineage <field-name> <field-value>",
		Aliases: []string{"l"},
		Short:   "Show all task lineages that match a certain field value",
		Long: `Show all task lineages that match a certain field value.
The field name is provided as first argument while the value is provided as second.
The default format used to print a task lineage is:
[comm] tid: ..., pid: ..., rptid: ..., rppid: ...
You can customize this format using the '--format' flag`,
		Example: `  - bpftree lineage tid 1 -> print lineage for tasks with tid=1
  - bpftree l t 1 -> print lineage for tasks with tid=1 (short form)
  - bpftree lineage comm systemd -> print lineage for tasks with comm=systemd
  - bpftree lineage comm systemd -f 't,p,r' -> print formatted tasks with comm=systemd`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			task.TasksLineage(args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(lineageCmd)
}
