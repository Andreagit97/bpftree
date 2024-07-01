package utils

import "fmt"

const (
	major uint64 = 0
	minor uint64 = 1
	patch uint64 = 0
)

// GetVersionString returns the current version of the tool as a string.
func GetVersionString() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// GetVersionMajor returns the Major version of the tool.
func GetVersionMajor() uint64 {
	return major
}

// GetVersionMinor returns the Minor version of the tool.
func GetVersionMinor() uint64 {
	return minor
}

// GetVersionPatch returns the Patch version of the tool.
func GetVersionPatch() uint64 {
	return patch
}

// GetVersion returns the current version of the tool.
func GetVersion() (uint64, uint64, uint64) {
	return GetVersionMajor(), GetVersionMinor(), GetVersionPatch()
}

// CheckCompatibleCapture checks if the capture file is compatible with the current version of the tool.
func CheckCompatibleCapture(captureMajor, captureMinor, capturePatch uint64) bool {
	return major == captureMajor && minor == captureMinor && patch >= capturePatch
}
