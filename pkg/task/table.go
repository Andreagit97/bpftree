package task

import (
	"fmt"
	"sort"

	"github.com/Andreagit97/bpftree/pkg/render"
)

var (
	tasksList []*task
)

func getTaskFromTid(tid int) *task {
	for _, task := range tasksList {
		if task.GetTid() == tid {
			return task
		}
	}
	return nil
}

func getSelectedTasksFromField(fieldType allowedField, value string) []*task {
	var selectedTasks []*task
	for _, task := range tasksList {
		if taskFieldString(task, fieldType) == value {
			selectedTasks = append(selectedTasks, task)
		}
	}
	return selectedTasks
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
			// It could be a missing thread or a thread outside the current namespace.
			// We want to report a warning only in the former case.
			if t.getVTid() != 1 && t.getVPid() != 1 {
				render.DisplayWarning(fmt.Sprintf("Parent with tid '%d' not found in the table during children computation", t.GetRealParentTid()))
			}
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
