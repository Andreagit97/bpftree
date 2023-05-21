package task

import (
	"bytes"
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