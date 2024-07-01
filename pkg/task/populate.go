package task

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/Andreagit97/bpftree/pkg/iterators"
	"github.com/Andreagit97/bpftree/pkg/render"
)

var (
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
			return err
		}
		tasksList = append(tasksList, &Task{
			Info:         taskInfo,
			Files:        make([]*FileInfo, 0),
			TaskChildren: make([]*Task, 0),
		})
	}

	reader.Close()

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
			return err
		}
		task := getTaskFromTid(int(fileInfo.TaskId))
		if task != nil {
			task.Files = append(task.Files, fileInfo)
		} else {
			render.DisplayWarning(fmt.Sprintf("Task with tid '%d' not found in the table during file info population.", fileInfo.TaskId))
		}
	}

	reader.Close()
	return nil
}

func populateTasksFromCaptureFile() error {
	// todo!: implement this function
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
			return int32(t.Files[i].GetFd()) < int32(t.Files[j].GetFd())
		})
	}

	return nil
}

func computeChildren() {
	/* we assign children according to the real parent tid */
	for _, t := range tasksList {
		/* init is the only one which doesn't have a parent */
		if t.GetRealParentTid() == 0 {
			continue
		}

		parent := getTaskFromTid(t.GetRealParentTid())
		if parent == nil {
			render.DisplayWarning(fmt.Sprintf("Task with tid '%d' not found in the table during children computation.", t.GetRealParentTid()))
			continue
		}
		parent.TaskChildren = append(parent.TaskChildren, t)
	}

	/* Order children according to their Tids */
	for _, t := range tasksList {
		if len(t.TaskChildren) == 0 {
			continue
		}
		sort.SliceStable(t.TaskChildren, func(i, j int) bool {
			return int32(t.TaskChildren[i].GetTid()) < int32(t.TaskChildren[j].GetTid())
		})
	}
}
