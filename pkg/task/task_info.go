package task

import (
	"encoding/binary"
	"io"
)

const (
	commLen    uint32 = 16
	exePathLen uint32 = 512
	cmdLineLen uint32 = 512
)

// TaskInfo is the struct sent by BPF side for each task.
type TaskInfo struct {
	/* These fields are mapped 1:1 to BPF side */
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

// func obtainTaskInfoField(reader io.ReadCloser, fieldSize uint32, data any) error {
// 	readTaskInfoSize += fieldSize
// 	if readTaskInfoSize != 0 && (readTaskInfoSize > totalTaskInfoSize) {
// 		return nil
// 	}
// 	return decodeByteEndianness(reader, fieldSize, data)
// }

// func parseTaskInfos(reader io.ReadCloser) error {
// 	for {
// 		taskInfo, err := parseTaskInfo(reader)
// 		if err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}
// 			return err
// 		}
// 		tasksList = append(tasksList, &task{
// 			Info: taskInfo,
// 		})
// 	}
// 	return nil
// }

func parseTaskInfo(reader io.ReadCloser) (*TaskInfo, error) {
	var t *TaskInfo

	// NativeEndian because we are reading on the same machine
	if err := binary.Read(reader, binary.NativeEndian, t); err != nil {
		return nil, err
	}
	return t, nil
}

// func parseTaskInfo(reader io.ReadCloser) (TaskInfo, error) {
// 	readTaskInfoSize = 0
// 	var t TaskInfo
// 	/* Tid */
// 	if err := obtainTaskInfoField(reader, 4, &t.Tid); err != nil {
// 		return t, err
// 	}

// 	/* Pid */
// 	if err := obtainTaskInfoField(reader, 4, &t.Pid); err != nil {
// 		return t, err
// 	}

// 	/* ParentTid */
// 	if err := obtainTaskInfoField(reader, 4, &t.ParentTid); err != nil {
// 		return t, err
// 	}

// 	/* ParentPid */
// 	if err := obtainTaskInfoField(reader, 4, &t.ParentPid); err != nil {
// 		return t, err
// 	}

// 	/* RealParentTid */
// 	if err := obtainTaskInfoField(reader, 4, &t.RealParentTid); err != nil {
// 		return t, err
// 	}

// 	/* RealParentPid */
// 	if err := obtainTaskInfoField(reader, 4, &t.RealParentPid); err != nil {
// 		return t, err
// 	}

// 	/* Comm */
// 	if err := obtainTaskInfoField(reader, commLen, &t.Comm); err != nil {
// 		return t, err
// 	}

// 	/* IsChildSubreaper */
// 	if err := obtainTaskInfoField(reader, 4, &t.IsChildSubreaper); err != nil {
// 		return t, err
// 	}

// 	/* Ns Level */
// 	if err := obtainTaskInfoField(reader, 4, &t.NsLevel); err != nil {
// 		return t, err
// 	}

// 	/* vtid */
// 	if err := obtainTaskInfoField(reader, 4, &t.VTid); err != nil {
// 		return t, err
// 	}

// 	/* vpid */
// 	if err := obtainTaskInfoField(reader, 4, &t.VPid); err != nil {
// 		return t, err
// 	}

// 	/* pgid */
// 	if err := obtainTaskInfoField(reader, 4, &t.Pgid); err != nil {
// 		return t, err
// 	}

// 	/* vpgid */
// 	if err := obtainTaskInfoField(reader, 4, &t.VPgid); err != nil {
// 		return t, err
// 	}

// 	/* sid */
// 	if err := obtainTaskInfoField(reader, 4, &t.Sid); err != nil {
// 		return t, err
// 	}

// 	/* vsid */
// 	if err := obtainTaskInfoField(reader, 4, &t.VSid); err != nil {
// 		return t, err
// 	}

// 	/* ExePath */
// 	if err := obtainTaskInfoField(reader, exePathLen, &t.ExePath); err != nil {
// 		return t, err
// 	}

// 	/* loginuid */
// 	if err := obtainTaskInfoField(reader, 8, &t.LoginUID); err != nil {
// 		return t, err
// 	}

// 	/* euid */
// 	if err := obtainTaskInfoField(reader, 8, &t.EUID); err != nil {
// 		return t, err
// 	}

// 	/* cmdline */
// 	if err := obtainTaskInfoField(reader, cmdLineLen, &t.CmdLine); err != nil {
// 		return t, err
// 	}

// 	/*
// 	 * Add all new fields here...
// 	 */
// 	return t, nil
// }
