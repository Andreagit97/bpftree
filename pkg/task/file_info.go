package task

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	filePathLen uint32 = 512
)

// FileInfo is the struct sent by BPF side for each task file.
type FileInfo struct {
	/* These fields are mapped 1:1 to BPF side */
	TaskId   int32
	Fd       uint32
	Ino      uint64
	MountId  uint32
	Flags    uint32
	FilePath [filePathLen]byte
}

func (f *FileInfo) GetFd() uint32 {
	return f.Fd
}

func (f *FileInfo) GetFilePath() string {
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

func (f *FileInfo) print() string {
	return fmt.Sprintf("(%d): path: %s, ino: %d, mount: %d, mount_id: %d", f.Fd, f.GetFilePath(), f.Ino, f.Flags, f.MountId)
}

func parseFileInfo(reader io.ReadCloser) (*FileInfo, error) {
	var f *FileInfo

	// NativeEndian because we are reading on the same machine
	if err := binary.Read(reader, binary.NativeEndian, f); err != nil {
		return nil, err
	}
	return f, nil
}
