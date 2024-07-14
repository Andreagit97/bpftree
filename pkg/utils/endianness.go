package utils

import "encoding/binary"

var (
	// CaptureEndianness all the capture files will be written in little endian.
	CaptureEndianness binary.ByteOrder = binary.LittleEndian
)
