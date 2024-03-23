package task

import (
	"fmt"
	"os"

	"github.com/enescakir/emoji"
)

type renderMode uint32

const (
	undefined renderMode = iota
	infoMode
	lineageMode
	treeMode
	formatMode
)

type allowedFields int

type fieldInfo struct {
	allowedNames []string
	description  string
	displayField func(t *task) string
	stringField  func(t *task) string
}

const (
	tidField allowedFields = iota
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

var allowedFieldsSlice = [maxField]fieldInfo{
	tidField: {
		allowedNames: []string{"tid", "t"},
		description:  "thread id (tid) of the current task (init namespace)",
		displayField: func(t *task) string {
			return fmt.Sprintf("t: %d, ", t.getTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getTid())
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
			return fmt.Sprintf("rpt: %d, ", t.getRealParentTid())
		},
		stringField: func(t *task) string {
			return fmt.Sprintf("%d", t.getRealParentTid())
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
		description:  "command line of the current task",
		displayField: func(t *task) string {
			return fmt.Sprintf("cmd: %s, ", t.getCmdLine())
		},
		stringField: func(t *task) string {
			return t.getCmdLine()
		},
	},
}

var (
	/* These are the emoji used in pretty print mode. */
	imageError     emoji.Emoji = emoji.RedCircle
	imageDownArrow emoji.Emoji = emoji.DownArrow
	imageReaper    emoji.Emoji = emoji.Skull
	imageNewspaper emoji.Emoji = emoji.RolledUpNewspaper
	imageTree      emoji.Emoji = emoji.PalmTree
	imageInfo      emoji.Emoji = emoji.Information
	imageLineage   emoji.Emoji = emoji.Scroll

	formatFields []string
	mode         = undefined
	displayGraph = func(a ...any) {
		fmt.Println(a...)
	}
	displayError = func(a ...any) {
		a = append([]any{imageError}, a...)
		fmt.Println(a...)
	}
)

func getFieldType(fieldName string) allowedFields {
	for fieldType, fieldInfo := range &allowedFieldsSlice {
		for _, allowedName := range fieldInfo.allowedNames {
			if allowedName == fieldName {
				return allowedFields(fieldType)
			}
		}
	}
	displayError(fmt.Sprintf("field '%s' is not allowed", fieldName))
	os.Exit(1)
	return maxField
}

func getFullFieldName(field string) string {
	fieldType := getFieldType(field)
	return allowedFieldsSlice[fieldType].allowedNames[0]
}

func DisablePrettyPrint() {
	imageError = "X"
	imageDownArrow = "V"
	imageReaper = "(R)"
	imageNewspaper = "*"
	imageTree = "-"
	imageInfo = "-"
	imageLineage = "-"
}

func setRenderMode(m renderMode) {
	/* We set the mode only if it is undefined
	 * otherwise there is the risk we overwrite the formatMode previously set
	 */
	if mode == undefined {
		mode = m
	}
}

// SetFormatFields is called to instruct bpftree on which fields it has to print.
// Moreover this method checks also if the provided fields are allowed or not.
func SetFormatFields(fields []string) {
	formatFields = fields
	for _, field := range formatFields {
		// if the field is not allowed `getFieldType` will log an error and terminate
		getFieldType(field)
	}
	setRenderMode(formatMode)
}

func (t *task) String() string {
	switch mode {
	case lineageMode, infoMode:
		return fmt.Sprintf("%s%v tid: %d, pid: %d, rptid: %d, rppid: %d",
			t.renderComm(),
			t.renderReaper(),
			t.getTid(),
			t.getPid(),
			t.getRealParentTid(),
			t.getRealParentPid())
	case treeMode:
		return fmt.Sprintf("%s%v tid: %d, pid: %d",
			t.renderComm(),
			t.renderReaper(),
			t.getTid(),
			t.getPid())
	case formatMode:
		finalLine := fmt.Sprintf("%s%v ", t.renderComm(), t.renderReaper())
		for _, field := range formatFields {
			finalLine += t.renderField(field)
		}
		/* truncate final space + , */
		finalLine = finalLine[:len(finalLine)-2]
		return finalLine
	default:
		displayError(fmt.Sprintf("Unknown render mode: '%d'", mode))
		os.Exit(1)
	}

	return ""
}

func (t *task) renderReaper() emoji.Emoji {
	if t.isChildSubReaper() {
		return imageReaper
	}
	/* if the task is not a child subreaper we don't add anything */
	return ""
}

/* main thread is rendered with "[]" while secondary threads are rendered with "{}". */
func (t *task) renderComm() string {
	if t.isMainThread() {
		return fmt.Sprintf("[%s]", t.getComm())
	}
	return fmt.Sprintf("{%s}", t.getComm())
}

func (t *task) renderField(field string) string {
	fieldType := getFieldType(field)
	return allowedFieldsSlice[int(fieldType)].displayField(t)
}
