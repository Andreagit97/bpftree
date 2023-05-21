package task

import (
	"io"
)

type Header struct {
	/* These fields are mapped 1:1 to BPF side */
	Version      uint32
	TaskInfoSize uint32
}

var (
	totalHeaderSize   uint32
	readHeaderSize    uint32
	totalTaskInfoSize uint32
)

func parseHeaderLen(reader io.ReadCloser) error {
	return decodeByteEndianness(reader, 4, &totalHeaderSize)
}

func obtainHeaderField(reader io.ReadCloser, fieldSize uint32, field any) error {
	readHeaderSize += fieldSize
	if readHeaderSize != 0 && (readHeaderSize > totalHeaderSize) {
		return nil
	}
	return decodeByteEndianness(reader, fieldSize, field)
}

func parseHeader(reader io.ReadCloser) error {
	var h Header
	readHeaderSize = 0
	/* Right now we don't need the header version we just read it */
	if err := obtainHeaderField(reader, 4, &h.Version); err != nil {
		return err
	}

	if err := obtainHeaderField(reader, 4, &h.TaskInfoSize); err != nil {
		return err
	}
	totalTaskInfoSize = h.TaskInfoSize

	/*
	 * Add all new fields here...
	 */

	return nil
}
