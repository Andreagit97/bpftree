package task

import (
	"bytes"
	"fmt"
	"os/user"
)

var (
	userList map[string]string
)

type task struct {
	Info     TaskInfo
	children []*task
}

func (t *task) isMainThread() bool {
	return t.Info.Pid == t.Info.Tid
}

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

func (t *task) isChildSubReaper() bool {
	return t.Info.IsChildSubreaper != 0
}

func (t *task) getParentTid() int {
	return int(t.Info.ParentTid)
}

func (t *task) getParentPid() int {
	return int(t.Info.ParentPid)
}

func (t *task) getRealParentTid() int {
	return int(t.Info.RealParentTid)
}

func (t *task) getRealParentPid() int {
	return int(t.Info.RealParentPid)
}

func (t *task) getTid() int {
	return int(t.Info.Tid)
}

func (t *task) getPid() int {
	return int(t.Info.Pid)
}

func (t *task) getChildren() []*task {
	return t.children
}

func (t *task) getNsLevel() int {
	return int(t.Info.NsLevel)
}

func (t *task) getVTid() int {
	return int(t.Info.VTid)
}

func (t *task) getVPid() int {
	return int(t.Info.VPid)
}

func (t *task) getPgid() int {
	return int(t.Info.Pgid)
}

func (t *task) getVPgid() int {
	return int(t.Info.VPgid)
}

func (t *task) getSid() int {
	return int(t.Info.Sid)
}

func (t *task) getVSid() int {
	return int(t.Info.VSid)
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

func (t *task) getLoginUIDName() string {
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

func (t *task) getEUIDName() string {
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
func (t *task) getCmdLine() string {
	// remove all the extra bytes at the end, if we have them
	fullCmdLine := bytes.Split(t.Info.CmdLine[:], []byte{0, 0})
	// we obtain a slice of slices with single arguments
	singleArgs := bytes.Split(fullCmdLine[0], []byte{0})
	joinedArgs := bytes.Join(singleArgs, []byte(","))
	return string(joinedArgs)
}

func init() {
	userList = make(map[string]string)
}
