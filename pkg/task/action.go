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

func printActionOnTask(t *Task, a action) {
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
			rparentTid := t.GetRealParentTid()
			t = getTaskFromTid(rparentTid)
			if t == nil {
				render.DisplayWarning("cannot find the thread with tid '%d' during lineage reconstruction", rparentTid)
				break
			}
		}
		// Leave a white line between a lineage and another
		render.DisplayGraph()
	case treeAction:
		// tree.SprintHrn generates a white line so we don't need to add an extra one.
		render.DisplayGraph(tree.SprintHrn(t))
	case fileAction:
		// todo!: implement this
		render.DisplayGraph("Implement me!")
	default:
		// We should never enter here
		render.DisplayError(fmt.Sprintf("Unknown action: '%d'", a))
		os.Exit(1)
	}

}

func printActionHeader(fieldType allowedField, value string, a action) {
	switch a {
	case infoAction:
		render.DisplayGraph(render.GetImageInfo(), "Task infoAction for '%s: %s'", getFullFieldName(fieldType), value)
	case lineageAction:
		render.DisplayGraph(render.GetImageLineage(), "Task lineageAction for '%s: %s'", getFullFieldName(fieldType), value)
	case treeAction:
		render.DisplayGraph(render.GetImageTree(), "Task treeAction for '%s: %s'", getFullFieldName(fieldType), value)
	case fileAction:
		render.DisplayGraph(render.GetImageFiles(), "Task Files for '%s: %s'", getFullFieldName(fieldType), value)
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
		render.DisplayError("invalid field name: ", field)
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
		render.DisplayError("There are no tasks with '%s' equals to '%s' in the system", getFullFieldName(fieldType), value)
		os.Exit(1)
	}

	// Print header before all the tasks
	printActionHeader(fieldType, value, a)

	for _, task := range selectedTasks {
		printActionOnTask(task, a)
	}
}

func TasksInfo(field, value string) {
	taskAction(field, value, infoAction)
}

func TasksFiles(field, value string) {
	taskAction(field, value, fileAction)
}

func TasksLineage(field, value string) {
	taskAction(field, value, lineageAction)
}

func TasksTree(field, value string) {
	taskAction(field, value, treeAction)
}
