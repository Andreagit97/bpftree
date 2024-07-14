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

type task struct {
	Info         *taskInfo
	TaskChildren []*task
	Files        []*fileInfo
}

func (t *task) isMainThread() bool {
	return t.Info.Pid == t.Info.Tid
}

func (t *task) isChildSubReaper() bool {
	return t.Info.IsChildSubreaper != 0
}

// Getters
func (t *task) getComm() string {
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

// The main thread is rendered with "[]" while secondary threads are rendered with "{}".
func (t *task) getFormattedComm() string {
	if t.isMainThread() {
		return fmt.Sprintf("[%s]", t.getComm())
	}
	return fmt.Sprintf("{%s}", t.getComm())
}

func (t *task) getParentTid() int {
	return int(t.Info.ParentTid)
}

func (t *task) getParentPid() int {
	return int(t.Info.ParentPid)
}

func (t *task) GetRealParentTid() int {
	return int(t.Info.RealParentTid)
}

func (t *task) getRealParentPid() int {
	return int(t.Info.RealParentPid)
}

func (t *task) GetTid() int {
	return int(t.Info.Tid)
}

func (t *task) getPid() int {
	return int(t.Info.Pid)
}

func (t *task) getNsLevel() int {
	return int(t.Info.NsLevel)
}

func (t *task) getVTid() int {
	return int(t.Info.VTid)
}

func (t *task) getVPid() int {
	// We use our user space table for 2 reasons:
	// - access virtual pids in the kernel could become tricky
	// - we don't want to send extra fields in the kernel since we can recover the virtual tids in userspace.
	if t.isMainThread() {
		return int(t.Info.VTid)
	}

	groupLeader := getTaskFromTid(int(t.Info.Pid))
	if groupLeader == nil {
		return int(-1)
	}
	return int(groupLeader.Info.VTid)
}

func (t *task) getPgid() int {
	return int(t.Info.Pgid)
}

func (t *task) getVPgid() int {
	groupLeader := getTaskFromTid(int(t.Info.Pgid))
	if groupLeader == nil {
		return int(-1)
	}
	return int(groupLeader.Info.VTid)
}

func (t *task) getSid() int {
	return int(t.Info.Sid)
}

func (t *task) getVSid() int {
	groupLeader := getTaskFromTid(int(t.Info.Sid))
	if groupLeader == nil {
		return int(-1)
	}
	return int(groupLeader.Info.VTid)
}

func (t *task) getExePath() string {
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

func (t *task) getLoginUID() int64 {
	return t.Info.LoginUID
}

func (t *task) getEUID() int64 {
	return t.Info.EUID
}

func searchUser(uid int64) string {
	if uid < 0 {
		return ""
	}

	// Search the username first in the map and then in the system
	UIDString := fmt.Sprintf("%d", uid)
	if username, ok := userList[UIDString]; ok {
		return username
	}

	u, err := user.LookupId(UIDString)
	if err == nil {
		userList[UIDString] = u.Username
		return u.Username
	}

	return ""
}

func (t *task) getLoginUIDName() string {
	return searchUser(t.Info.LoginUID)
}

func (t *task) getEUIDName() string {
	return searchUser(t.Info.EUID)
}

// todo!: this is simple scratch implementation it doesn't cover all the cases.
// Moreover we need to address the `set_proctitle` call in the kernel code.
func (t *task) getCmdLine() string {
	// remove all the extra bytes at the end, if we have them
	fullCmdLine := bytes.Split(t.Info.CmdLine[:], []byte{0, 0})
	// we obtain a slice of slices with single arguments
	singleArgs := bytes.Split(fullCmdLine[0], []byte{0})
	joinedArgs := bytes.Join(singleArgs, []byte(","))
	return string(joinedArgs)
}

// Here we implement tree.Node interface with `Children` and `Data` methods
func (t *task) Children() []tree.Node {
	var treeNodes []tree.Node
	for _, child := range t.TaskChildren {
		treeNodes = append(treeNodes, tree.Node(child))
	}
	return treeNodes
}

func (t *task) Data() interface{} {
	return t.String()
}

func (t *task) getReaperImage() emoji.Emoji {
	if t.isChildSubReaper() {
		return render.GetImageReaper()
	}
	/* if the task is not a child subreaper we don't add anything */
	return ""
}

func init() {
	userList = make(map[string]string)
}
