package task

import (
	"bytes"
	"fmt"
	"os/user"

	"github.com/Andreagit97/bpftree/pkg/render"
	"github.com/enescakir/emoji"
	"github.com/shivamMg/ppds/tree"
)

var (
	userList map[string]string
)

type Task struct {
	Info         *TaskInfo
	TaskChildren []*Task
	Files        []*FileInfo
}

func (t *Task) isMainThread() bool {
	return t.Info.Pid == t.Info.Tid
}

func (t *Task) getComm() string {
	/* We truncate the comm after the \0 to avoid printing
	 * useless end characters
	 */
	index := bytes.IndexByte(t.Info.Comm[:], byte(0))
	if index == -1 {
		/* There is no terminator */
		return string(t.Info.Comm[:])
	}
	return string(t.Info.Comm[:index])
}

func (t *Task) isChildSubReaper() bool {
	return t.Info.IsChildSubreaper != 0
}

func (t *Task) getParentTid() int {
	return int(t.Info.ParentTid)
}

func (t *Task) getParentPid() int {
	return int(t.Info.ParentPid)
}

func (t *Task) GetRealParentTid() int {
	return int(t.Info.RealParentTid)
}

func (t *Task) getRealParentPid() int {
	return int(t.Info.RealParentPid)
}

func (t *Task) GetTid() int {
	return int(t.Info.Tid)
}

func (t *Task) getPid() int {
	return int(t.Info.Pid)
}

func (t *Task) getChildren() []*Task {
	return t.TaskChildren
}

func (t *Task) getNsLevel() int {
	return int(t.Info.NsLevel)
}

func (t *Task) getVTid() int {
	return int(t.Info.VTid)
}

func (t *Task) getVPid() int {
	// todo!: we need to implement this

	// if t.isMainThread() {
	// 	return int(t.Info.VTid)
	// } else {
	// 	groupLeader := getTaskFromTid(..., t.Info.Pid)
	// 	return int(groupLeader.Info.VTid)
	// }

	return int(0)
}

func (t *Task) getPgid() int {
	return int(t.Info.Pgid)
}

func (t *Task) getVPgid() int {

	// todo!: we need to implement this

	// 	groupLeader := getTaskFromTid(..., t.Info.Pgid)
	// 	return int(groupLeader.Info.VTid)

	return int(0)
}

func (t *Task) getSid() int {
	return int(t.Info.Sid)
}

func (t *Task) getVSid() int {
	// todo!: we need to implement this

	// 	groupLeader := getTaskFromTid(..., t.Info.Sid)
	// 	return int(groupLeader.Info.VTid)

	return int(0)
}

func (t *Task) getExePath() string {
	/* We truncate the exePath after the \0 to avoid printing
	 * useless end characters
	 */
	index := bytes.IndexByte(t.Info.ExePath[:], byte(0))
	if index == -1 {
		/* There is no terminator */
		return string(t.Info.ExePath[:])
	}
	return string(t.Info.ExePath[:index])
}

func (t *Task) getLoginUID() int64 {
	return t.Info.LoginUID
}

func (t *Task) getEUID() int64 {
	return t.Info.EUID
}

func (t *Task) getLoginUIDName() string {
	/* Search the username first in the map and then in the system */
	if t.Info.LoginUID < 0 {
		return ""
	}

	loginUIDString := fmt.Sprintf("%d", t.Info.LoginUID)
	if username, ok := userList[loginUIDString]; ok {
		return username
	}

	u, err := user.LookupId(loginUIDString)
	if err == nil {
		userList[loginUIDString] = u.Username
		return u.Username
	}

	return ""
}

func (t *Task) getEUIDName() string {
	/* Search the username first in the map and then in the system */
	if t.Info.EUID < 0 {
		return ""
	}

	eUIDString := fmt.Sprintf("%d", t.Info.EUID)

	if username, ok := userList[eUIDString]; ok {
		return username
	}

	u, err := user.LookupId(eUIDString)
	if err == nil {
		userList[eUIDString] = u.Username
		return u.Username
	}

	return ""
}

// todo!: this is simple scratch implementation it doesn't cover all the cases.
// Moreover we need to address the `set_proctitle` call in the kernel code.
func (t *Task) getCmdLine() string {
	// remove all the extra bytes at the end, if we have them
	fullCmdLine := bytes.Split(t.Info.CmdLine[:], []byte{0, 0})
	// we obtain a slice of slices with single arguments
	singleArgs := bytes.Split(fullCmdLine[0], []byte{0})
	joinedArgs := bytes.Join(singleArgs, []byte(","))
	return string(joinedArgs)
}

// Here we implement tree.Node interface with `Children` and `Data` methods
func (t *Task) Children() []tree.Node {
	var treeNodes []tree.Node
	for _, child := range t.TaskChildren {
		treeNodes = append(treeNodes, tree.Node(child))
	}
	return treeNodes
}

func (t *Task) Data() interface{} {
	return t.String()
}

func (t *Task) getReaperImage() emoji.Emoji {
	if t.isChildSubReaper() {
		return render.GetImageReaper()
	}
	/* if the task is not a child subreaper we don't add anything */
	return ""
}

// The main thread is rendered with "[]" while secondary threads are rendered with "{}".
func (t *Task) getFormattedComm() string {
	if t.isMainThread() {
		return fmt.Sprintf("[%s]", t.getComm())
	}
	return fmt.Sprintf("{%s}", t.getComm())
}

func init() {
	userList = make(map[string]string)
}
