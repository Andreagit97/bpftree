package cmd

import (
	"github.com/Andreagit97/bpftree/pkg/task"
	"github.com/spf13/cobra"
)

var (
	treeCmd = &cobra.Command{
		Use:     "tree <field-name> <field-value>",
		Aliases: []string{"t"},
		Short:   "Show all task trees that match a certain field value",
		Long: `Show all task trees that match a certain field value.
The field name is provided as first argument while the value is provided as second.
The default format used to print a task info is:
[comm] tid: ..., pid: ...
You can customize this format using the '--format' flag`,
		Example: `  - bpftree tree tid 1 -> print trees for tasks with tid=1
  - bpftree t t 1 -> print trees for tasks with tid=1 (short form)
  - bpftree tree comm systemd -> print trees for tasks with comm=systemd
  - bpftree tree comm systemd -f 't,p,r' -> print formatted tasks with comm=systemd`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			task.TasksTree(args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(treeCmd)
}
