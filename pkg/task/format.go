package task

import (
	"fmt"
	"os"

	"github.com/Andreagit97/bpftree/pkg/render"
)

var (
	// FormatFields is the list of fields to be printed in the format mode
	FormatFields      []string
	formatFieldsTypes []allowedField
	printTaskCb       = func(t *task) string {
		render.DisplayError("printTaskCb not set")
		os.Exit(1)
		return ""
	}
)

func setFormatMode(m action) {
	if len(FormatFields) != 0 {
		if m == fileAction {
			render.DisplayError("format fields are not yet supported for 'files,fds' command")
			os.Exit(1)
		}

		// We validate the fields
		for _, field := range FormatFields {
			fieldType := getFieldTypeFromName(field)
			if isUnknownField(fieldType) {
				render.DisplayError(fmt.Sprintf("invalid field name in the format fields '%s'", field))
				os.Exit(1)
			}
			formatFieldsTypes = append(formatFieldsTypes, fieldType)
		}

		// todo!: probably we can improve this callback without using a for loop
		printTaskCb = func(t *task) string {
			finalLine := fmt.Sprintf("%s%v ", t.getFormattedComm(), t.getReaperImage())
			for _, fieldType := range formatFieldsTypes {
				finalLine += taskFieldDisplay(t, fieldType)
			}
			/* truncate final space + , */
			finalLine = finalLine[:len(finalLine)-2]
			return finalLine
		}
		return
	}

	switch m {
	case infoAction, lineageAction:
		printTaskCb = func(t *task) string {
			return fmt.Sprintf("%s%v tid: %d, pid: %d, rptid: %d, rppid: %d",
				t.getFormattedComm(),
				t.getReaperImage(),
				t.GetTid(),
				t.getPid(),
				t.GetRealParentTid(),
				t.getRealParentPid())
		}
	case treeAction:
		printTaskCb = func(t *task) string {
			return fmt.Sprintf("%s%v tid: %d, pid: %d",
				t.getFormattedComm(),
				t.getReaperImage(),
				t.GetTid(),
				t.getPid())
		}
	case fileAction:
		// we won't print the task so we don't need to set the callback
		return
	}
}

func (t *task) String() string {
	return printTaskCb(t)
}
