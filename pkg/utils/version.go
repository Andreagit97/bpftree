package utils

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Version is the struct that represents the version of the tool.
type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

var toolVersion = Version{
	Major: 0,
	Minor: 1,
	Patch: 0,
}

func getVersionString(v Version) string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// GetToolVersion returns the current version of the tool as a string.
func GetToolVersion() string {
	return getVersionString(toolVersion)
}

// checkCompatibleCapture checks if the capture file is compatible with the current version of the tool.
// Only a tool with the same Major and Minor and with a greater or equal Patch can read the file.
func checkCompatibleCapture(captureVersion Version) bool {
	return toolVersion.Major == captureVersion.Major &&
		toolVersion.Minor == captureVersion.Minor &&
		toolVersion.Patch >= captureVersion.Patch
}

// DumpToolVersion writes the current version of the tool to the file.
func DumpToolVersion(f *os.File) error {
	if err := binary.Write(f, CaptureEndianness, &toolVersion); err != nil {
		return fmt.Errorf("cannot dump the version into the file. %w", err)
	}
	return nil
}

// RetrieveAndCheckCompatibleVersion reads the version of the capture file and checks if it is compatible with the current version of the tool.
func RetrieveAndCheckCompatibleVersion(f *os.File) error {
	var captureVersion Version
	if err := binary.Read(f, CaptureEndianness, &captureVersion); err != nil {
		return fmt.Errorf("cannot read the version of the capture file. %w", err)
	}

	if !checkCompatibleCapture(captureVersion) {
		return fmt.Errorf("the capture file is not compatible with the current version of the tool. "+
			"Tool version: %s, Capture version: %s", GetToolVersion(), getVersionString(captureVersion))
	}

	return nil
}
