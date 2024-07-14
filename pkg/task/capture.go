package task

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Andreagit97/bpftree/pkg/iterators"
	"github.com/Andreagit97/bpftree/pkg/render"
	"github.com/Andreagit97/bpftree/pkg/utils"
)

// We have 3 uint64 for the version
const versionOffset = 3 * 8

/* Note: today we dump all the files and all the tasks but maybe in the future we could dump only some of them.
 *
 * Update this if the schema changes! Please note that the version should be always dumped
 *
 * Format we dump the capture file (v0.2.0):
 * --------------------------------
 * - Major version (uint64)
 * - Minor version (uint64)
 * - Patch version (uint64)
 * --------------------------------
 * - Number of task info structs (uint64)
 * - Number of file info structs (uint64)
 * - All the task info structs
 * - All the file info structs
 */

// DumpToFile dumps all the task and file info into a file
func DumpToFile(path string) error {

	// Open the file to write
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to get absolute path. %v", err)
	}

	f, err := os.Create(filepath.Clean(absPath))
	if err != nil {
		return fmt.Errorf("cannot create dump file. %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 1. Write the version of the tool into the file.
	if err := utils.DumpToolVersion(f); err != nil {
		return err
	}

	// 2. Write the number of task info into the file
	// Initially we don't know how many info we have, so we write 0
	numTasks := uint64(0)
	if err := binary.Write(f, utils.CaptureEndianness, &numTasks); err != nil {
		return fmt.Errorf("cannot write the initial 0 for task info. %w", err)
	}

	// 3. Write the number of file info into the file
	// Initially we don't know how many info we have, so we write 0
	numFiles := uint64(0)
	if err := binary.Write(f, utils.CaptureEndianness, &numFiles); err != nil {
		return fmt.Errorf("cannot write the initial 0 for file info. %w", err)
	}

	// 4. Write all the task info into the file
	reader, err := iterators.GetTasksReader()
	if err != nil {
		return err
	}

	for {
		taskInfo, err := parseTaskInfo(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("cannot parse task info. %w", err)
		}
		numTasks++
		if err := taskInfo.dumpIntoCapture(f); err != nil {
			return err
		}
	}
	if err := reader.Close(); err != nil {
		return fmt.Errorf("cannot close the reader for task info. %w", err)
	}

	// 5. Write all the file info into the file
	reader, err = iterators.GetFilesReader()
	if err != nil {
		return err
	}

	for {
		fileInfo, err := parseFileInfo(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("cannot parse file info. %w", err)
		}
		numFiles++
		if err := fileInfo.dumpIntoCapture(f); err != nil {
			return err
		}
	}

	// 6. Seek after the version and write the number of task info and file info
	if _, err := f.Seek(versionOffset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	if err := binary.Write(f, utils.CaptureEndianness, &numTasks); err != nil {
		return fmt.Errorf("cannot write the number of task info. %w", err)
	}

	if err := binary.Write(f, utils.CaptureEndianness, &numFiles); err != nil {
		return fmt.Errorf("cannot write the number of file info. %w", err)
	}

	if err := reader.Close(); err != nil {
		return fmt.Errorf("cannot close the reader for file info. %w", err)
	}
	render.DisplayGraph(render.GetImageNewspaper(), "Capture correctly dumped:", filepath.Clean(absPath))
	return nil
}
