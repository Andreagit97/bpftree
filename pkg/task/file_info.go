package task

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Andreagit97/bpftree/pkg/utils"
)

const (
	filePathLen uint32 = 512
)

const (
	// Mountable is the flag to check if the file is mountable
	Mountable = 1 << iota // 1 << 0 which is 00000001
	// UpperLayer is the flag to check if the file is in the upper layer
	UpperLayer // 1 << 1 which is 00000010
	// LowerLayer is the flag to check if the file is in the lower layer
	LowerLayer // 1 << 2 which is 00000100
)

type fileFlags uint32

func (f *fileInfo) convertFlagToString() string {
	var finalLine string

	if f.Flags&Mountable != 0 {
		finalLine += "MOUNTABLE, "
	} else {
		finalLine += "NOT_MOUNTABLE, "
	}

	if f.Flags&UpperLayer != 0 {
		finalLine += "UPPER_LAYER, "
	}

	if f.Flags&LowerLayer != 0 {
		finalLine += "LOWER_LAYER, "
	}

	// truncate final space + ,
	finalLine = finalLine[:len(finalLine)-2]
	return finalLine
}

/* These fields are mapped 1:1 to BPF side */

// fileInfo is the struct sent by BPF side for each task file.
type fileInfo struct {
	TaskID   int32
	Fd       uint32
	Ino      uint64
	MountID  uint32
	Flags    fileFlags
	FilePath [filePathLen]byte
}

func (f *fileInfo) getFd() uint32 {
	return f.Fd
}

func (f *fileInfo) getFilePath() string {
	/* We truncate the exePath after the \0 to avoid printing
	 * useless end characters
	 */
	index := bytes.IndexByte(f.FilePath[:], byte(0))
	if index == -1 {
		/* There is no terminator */
		return string(f.FilePath[:])
	}
	return string(f.FilePath[:index])
}

func (f *fileInfo) String() string {
	return fmt.Sprintf("(%d): path: %s, ino: %d, flags: (%s), mount_id: %d", f.Fd, f.getFilePath(), f.Ino, f.convertFlagToString(), f.MountID)
}

func parseFileInfo(reader io.ReadCloser) (*fileInfo, error) {
	var f fileInfo

	// NativeEndian because we are reading on the same machine
	if err := binary.Read(reader, binary.NativeEndian, &f); err != nil {
		return nil, fmt.Errorf("cannot read file info. %w", err)
	}
	return &f, nil
}

func (f *fileInfo) dumpIntoCapture(file *os.File) error {
	if err := binary.Write(file, utils.CaptureEndianness, f); err != nil {
		return fmt.Errorf("cannot write file info. %w", err)
	}
	return nil
}
