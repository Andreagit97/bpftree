package task

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Andreagit97/bpftree/pkg/utils"
)

const (
	commLen    uint32 = 16
	exePathLen uint32 = 512
	cmdLineLen uint32 = 512
)

/* These fields are mapped 1:1 to BPF side */

// taskInfo is the struct sent by BPF side for each task.
type taskInfo struct {
	Tid              int32
	Pid              int32
	Pgid             int32
	Sid              int32
	VTid             int32
	ParentTid        int32
	ParentPid        int32
	RealParentTid    int32
	RealParentPid    int32
	IsChildSubreaper uint8
	NsLevel          uint32
	LoginUID         int64
	EUID             int64
	Comm             [commLen]byte
	ExePath          [exePathLen]byte
	CmdLine          [cmdLineLen]byte
}

func parseTaskInfo(reader io.ReadCloser) (*taskInfo, error) {
	var t taskInfo

	// NativeEndian because we are reading on the same machine
	if err := binary.Read(reader, binary.NativeEndian, &t); err != nil {
		return nil, fmt.Errorf("cannot read task info. %w", err)
	}
	return &t, nil
}

func (t *taskInfo) dumpIntoCapture(file *os.File) error {
	if err := binary.Write(file, utils.CaptureEndianness, t); err != nil {
		return fmt.Errorf("cannot write task info. %w", err)
	}
	return nil
}
