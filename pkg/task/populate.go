package task

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Andreagit97/bpftree/pkg/iterators"
	"github.com/Andreagit97/bpftree/pkg/render"
	"github.com/Andreagit97/bpftree/pkg/utils"
)

var (
	// CaptureFilePath is the path to be used to replay the capture
	CaptureFilePath = ""
)

func populateTasksFromIterators() error {
	// Parse the task infos
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
		tasksList = append(tasksList, &task{
			Info:         taskInfo,
			Files:        make([]*fileInfo, 0),
			TaskChildren: make([]*task, 0),
		})
	}
	if err := reader.Close(); err != nil {
		return fmt.Errorf("cannot close the reader for task info. %w", err)
	}

	// Parse the file infos
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
		task := getTaskFromTid(int(fileInfo.TaskID))
		// it is possible that some files belong to new tasks that we don't have in the table, we skip them
		if task != nil {
			task.Files = append(task.Files, fileInfo)
		}
	}

	if err := reader.Close(); err != nil {
		return fmt.Errorf("cannot close the reader for file info. %w", err)
	}
	return nil
}

func populateTasksFromCaptureFile() error {
	file, err := os.Open(CaptureFilePath)
	if err != nil {
		return fmt.Errorf("unable to open file '%s'. %v", CaptureFilePath, err)
	}

	if err := utils.RetrieveAndCheckCompatibleVersion(file); err != nil {
		return err
	}

	numTasks := uint64(0)
	numFiles := uint64(0)
	if err := binary.Read(file, utils.CaptureEndianness, &numTasks); err != nil {
		return fmt.Errorf("cannot read the number of task info. %w", err)
	}
	if err := binary.Read(file, utils.CaptureEndianness, &numFiles); err != nil {
		return fmt.Errorf("cannot read the number of file info. %w", err)
	}

	for i := uint64(0); i < numTasks; i++ {
		taskInfo, err := parseTaskInfo(file)
		if err != nil {
			return fmt.Errorf("cannot parse task info. %w", err)
		}
		tasksList = append(tasksList, &task{
			Info:         taskInfo,
			Files:        make([]*fileInfo, 0),
			TaskChildren: make([]*task, 0),
		})
	}

	for i := uint64(0); i < numFiles; i++ {
		fileInfo, err := parseFileInfo(file)
		if err != nil {
			return fmt.Errorf("cannot parse file info. %w", err)
		}
		task := getTaskFromTid(int(fileInfo.TaskID))
		if task != nil {
			task.Files = append(task.Files, fileInfo)
		} else {
			render.DisplayWarning(fmt.Sprintf("Task with tid '%d' not found in the table during file info population.", fileInfo.TaskID))
		}
	}
	return nil
}

func populateTasks() error {
	var err error

	if CaptureFilePath == "" {
		/* We need to read from the system with BPF. */
		err = populateTasksFromIterators()
	} else {
		/* We can use the provided file. */
		err = populateTasksFromCaptureFile()
	}

	if err != nil {
		return err
	}

	// Order the list of tasks by tid
	sort.SliceStable(tasksList, func(i, j int) bool {
		return tasksList[i].GetTid() < tasksList[j].GetTid()
	})

	// Compute the children of each task
	computeChildren()

	// Order files by fd
	for _, t := range tasksList {
		if len(t.Files) == 0 {
			continue
		}
		sort.SliceStable(t.Files, func(i, j int) bool {
			return int32(t.Files[i].getFd()) < int32(t.Files[j].getFd())
		})
	}

	return nil
}
