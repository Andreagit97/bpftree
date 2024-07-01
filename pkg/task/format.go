package task

import (
	"fmt"
	"os"

	"github.com/Andreagit97/bpftree/pkg/render"
)

var (
	FormatFields      []string
	formatFieldsTypes []allowedField
	printTaskCb       = func(t *Task) string {
		render.DisplayError("printTaskCb not set")
		os.Exit(1)
		return ""
	}
)

func setFormatMode(m action) {
	if len(FormatFields) != 0 {
		// We validate the fields
		for _, field := range FormatFields {
			fieldType := getFieldTypeFromName(field)
			if isUnknownField(fieldType) {
				render.DisplayError("invalid field name in the format fields: ", field)
				os.Exit(1)
			}
			formatFieldsTypes = append(formatFieldsTypes, fieldType)
		}

		// todo!: probably we can improve this callback without using a for loop
		printTaskCb = func(t *Task) string {
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
	case infoAction:
	case lineageAction:
		printTaskCb = func(t *Task) string {
			return fmt.Sprintf("%s%v tid: %d, pid: %d, rptid: %d, rppid: %d",
				t.getFormattedComm(),
				t.getReaperImage(),
				t.GetTid(),
				t.getPid(),
				t.GetRealParentTid(),
				t.getRealParentPid())
		}
	case treeAction:
		printTaskCb = func(t *Task) string {
			return fmt.Sprintf("%s%v tid: %d, pid: %d",
				t.getFormattedComm(),
				t.getReaperImage(),
				t.GetTid(),
				t.getPid())
		}
	case fileAction:
		printTaskCb = func(t *Task) string {
			// todo!: implement this
			return ""
		}
	}
}

func (t *Task) String() string {
	return printTaskCb(t)
}
