package task

import (
	"fmt"
	"strings"

	"github.com/Andreagit97/bpftree/pkg/render"
)

type allowedField int

// These are task fields, we should have also file fields
const (
	unknownField allowedField = iota
	tidField
	vTidField
	pidField
	vPidField
	pgidField
	vPgidField
	sidField
	vSidField
	parentTidField
	parentPidField
	realParentTidField
	realParentPidField
	commField
	reaperField
	nsLevelField
	exePathField
	loginUIDField
	eUIDField
	cmdLineField
	maxField
)

type fieldInfo struct {
	allowedNames []string
	description  string
	displayField func(t *task) string
	stringField  func(t *task) string
}

var allowedFieldsSlice = [maxField]fieldInfo{
	tidField: {
		allowedNames: []string{"tid", "t"},
		description:  "thread id (tid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("t: %d, ", t.GetTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.GetTid())
		},
	},

	vTidField: {
		allowedNames: []string{"vtid", "vt"},
		description:  "thread id (tid) of the current task (task namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("vt: %d, ", t.getVTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getVTid())
		},
	},

	pidField: {
		allowedNames: []string{"pid", "p"},
		description:  "process id (pid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("p: %d, ", t.getPid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getPid())
		},
	},

	vPidField: {
		allowedNames: []string{"vpid", "vp"},
		description:  "process id (pid) of the current task (task namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("vp: %d, ", t.getVPid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getVPid())
		},
	},

	pgidField: {
		allowedNames: []string{"pgid", "pg"},
		description:  "process group id (pgid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("pg: %d, ", t.getPgid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getPgid())
		},
	},

	vPgidField: {
		allowedNames: []string{"vpgid", "vpg"},
		description:  "process group id (pgid) of the current task (task namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("vpg: %d, ", t.getVPgid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getVPgid())
		},
	},

	sidField: {
		allowedNames: []string{"sid", "s"},
		description:  "session id (sid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("s: %d, ", t.getSid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getSid())
		},
	},

	vSidField: {
		allowedNames: []string{"vsid", "vs"},
		description:  "session id (sid) of the current task (task namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("vs: %d, ", t.getVSid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getVSid())
		},
	},

	parentTidField: {
		allowedNames: []string{"ptid", "pt"},
		description:  "parent thread id (ptid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("pt: %d, ", t.getParentTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getParentTid())
		},
	},

	parentPidField: {
		allowedNames: []string{"ppid", "pp"},
		description:  "parent process id (ppid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("pp: %d, ", t.getParentPid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getParentPid())
		},
	},

	realParentTidField: {
		allowedNames: []string{"rptid", "rpt"},
		description:  "real parent thread id (rptid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("rpt: %d, ", t.GetRealParentTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.GetRealParentTid())
		},
	},

	realParentPidField: {
		allowedNames: []string{"rppid", "rpp"},
		description:  "real parent process id (rppid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("rpp: %d, ", t.getRealParentPid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getRealParentPid())
		},
	},

	commField: {
		allowedNames: []string{"comm", "c"},
		description:  "human readeable process name ('task->comm')",
		displayField: func(t *task) string {
			return fmt.Sprintf("c: %s, ", t.getComm())
		},
		stringField: func(t *task) string {
			return t.getComm()
		},
	},

	reaperField: {
		allowedNames: []string{"reaper", "r"},
		description:  "true if the current process is a child_sub_reaper",
		displayField: func(t *task) string {
			return fmt.Sprintf("r: %v, ", t.isChildSubReaper())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%v", t.isChildSubReaper())
		},
	},

	nsLevelField: {
		allowedNames: []string{"ns_level", "ns"},
		description:  "pid namespace level of the actual thread",
		displayField: func(t *task) string {
			return fmt.Sprintf("ns: %d, ", t.getNsLevel())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getNsLevel())
		},
	},

	exePathField: {
		allowedNames: []string{"exepath", "e"},
		description:  "full executable path of the current task",
		displayField: func(t *task) string {
			return fmt.Sprintf("e: %s, ", t.getExePath())
		},
		stringField: func(t *task) string {
			return t.getExePath()
		},
	},

	loginUIDField: {
		allowedNames: []string{"loginuid", "lu"},
		description:  "UID of the user that interacted with a login service",
		displayField: func(t *task) string {
			return fmt.Sprintf("lu: %d(%s), ", t.getLoginUID(), t.getLoginUIDName())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getLoginUID())
		},
	},

	eUIDField: {
		allowedNames: []string{"euid", "eu"},
		description:  "Effective UID",
		displayField: func(t *task) string {
			return fmt.Sprintf("eu: %d(%s), ", t.getEUID(), t.getEUIDName())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getEUID())
		},
	},

	cmdLineField: {
		allowedNames: []string{"cmdline", "cmd"},
		description:  "Command line of the current task",
		displayField: func(t *task) string {
			return fmt.Sprintf("cmd: %s, ", t.getCmdLine())
		},
		stringField: func(t *task) string {
			return t.getCmdLine()
		},
	},
}

// this should be the unique method called to validate the field names
func getFieldTypeFromName(fieldName string) allowedField {
	for fieldType, fieldInfo := range &allowedFieldsSlice {
		for _, allowedName := range fieldInfo.allowedNames {
			if allowedName == fieldName {
				return allowedField(fieldType)
			}
		}
	}
	return unknownField
}

func taskFieldDisplay(t *task, fieldType allowedField) string {
	return allowedFieldsSlice[fieldType].displayField(t)
}

func taskFieldString(t *task, fieldType allowedField) string {
	return allowedFieldsSlice[fieldType].stringField(t)
}

func getFullFieldName(fieldType allowedField) string {
	return allowedFieldsSlice[fieldType].allowedNames[0]
}

func isUnknownField(fieldType allowedField) bool {
	return fieldType == unknownField
}

// RenderAllowedFields prints the allowed fields in a markdown table
func RenderAllowedFields() {
	render.DisplayGraph("| Fields | Description |")
	render.DisplayGraph("| ------ | ----------- |")
	for i, info := range &allowedFieldsSlice {
		// We skip the first undefined field
		if i == 0 {
			continue
		}
		render.DisplayGraph("|", strings.Join(info.allowedNames, ","), "|", info.description, "|")
	}
}
