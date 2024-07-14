package task

import (
	"fmt"
	"os"

	"github.com/Andreagit97/bpftree/pkg/render"
	"github.com/shivamMg/ppds/tree"
)

type action uint32

const (
	infoAction action = iota
	lineageAction
	treeAction
	fileAction
)

func printActionOnTask(t *task, a action) {
	switch a {
	case infoAction:
		render.DisplayGraph(render.GetImageNewspaper(), t)
	case lineageAction:
		for {
			/* Print current task info */
			render.DisplayGraph(render.GetImageDownArrow(), t)

			/* We reached the last ancestor, stop */
			if t.GetRealParentTid() == 0 {
				break
			}

			/* go to the real parent */
			current := t
			t = getTaskFromTid(current.GetRealParentTid())
			if t == nil {
				// It could be a missing thread or a thread outside the current namespace.
				// We want to report a warning only in the former case.
				if current.getVTid() != 1 && current.getVPid() != 1 {
					render.DisplayWarning(fmt.Sprintf("cannot find the parent with tid '%d' during lineage reconstruction", current.GetRealParentTid()))
				}
				break
			}
		}
		// Leave a white line between a lineage and another
		render.DisplayGraph()
	case treeAction:
		// tree.SprintHrn generates a white line so we don't need to add an extra one.
		render.DisplayGraph(tree.SprintHrn(t))
	case fileAction:
		render.DisplayGraph(render.GetImageFile(), fmt.Sprintf("Files for task %s(%d-%d):", t.getFormattedComm(), t.GetTid(), t.getPid()))
		if len(t.Files) == 0 {
			render.DisplayGraph("No files found for this task.")
		} else {
			for _, f := range t.Files {
				render.DisplayGraph(f)
			}
		}
		// Leave a white line between a task and another
		render.DisplayGraph()
	default:
		// We should never enter here
		render.DisplayError("Unknown action: ", a)
		os.Exit(1)
	}

}

func printActionHeader(fieldType allowedField, value string, a action) {
	switch a {
	case infoAction:
		render.DisplayGraph(render.GetImageInfo(), fmt.Sprintf("Task Info for '%s=%s'", getFullFieldName(fieldType), value))
	case lineageAction:
		render.DisplayGraph(render.GetImageLineage(), fmt.Sprintf("Task Lineage for '%s=%s'", getFullFieldName(fieldType), value))
	case treeAction:
		render.DisplayGraph(render.GetImageTree(), fmt.Sprintf("Task Tree for '%s=%s'", getFullFieldName(fieldType), value))
	case fileAction:
		render.DisplayGraph(render.GetImageFolder(), fmt.Sprintf("Task Files for '%s=%s'", getFullFieldName(fieldType), value))
	default:
		// We should never enter here
		render.DisplayError(fmt.Sprintf("Unknown action: '%d'", a))
		os.Exit(1)
	}
}

func taskAction(field, value string, a action) {
	// we disable the pretty print mode if the user has requested it
	render.ConfigureRendering()

	// The validation of the input field is done here
	fieldType := getFieldTypeFromName(field)
	if isUnknownField(fieldType) {
		render.DisplayError(fmt.Sprintf("invalid field name '%s'", field))
		os.Exit(1)
	}

	setFormatMode(a)

	// Populate all tasks
	err := populateTasks()
	if err != nil {
		render.DisplayError(err)
		os.Exit(1)
	}

	// Select only some of them
	selectedTasks := getSelectedTasksFromField(fieldType, value)
	if len(selectedTasks) == 0 {
		render.DisplayError(fmt.Sprintf("There are no tasks with '%s=%s' in the system", getFullFieldName(fieldType), value))
		os.Exit(1)
	}

	// Print header before all the tasks
	printActionHeader(fieldType, value, a)

	for _, task := range selectedTasks {
		printActionOnTask(task, a)
	}
}

// TasksInfo prints the task info for the given field and value.
func TasksInfo(field, value string) {
	taskAction(field, value, infoAction)
}

// TasksFiles prints the task files for the given field and value.
func TasksFiles(field, value string) {
	taskAction(field, value, fileAction)
}

// TasksLineage prints the task lineage for the given field and value.
func TasksLineage(field, value string) {
	taskAction(field, value, lineageAction)
}

// TasksTree prints the task tree for the given field and value.
func TasksTree(field, value string) {
	taskAction(field, value, treeAction)
}
