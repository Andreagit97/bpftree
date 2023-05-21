package task

import (
	"fmt"

	"github.com/shivamMg/ppds/tree"
)

// Here we implement tree.Node interface with `Children` and `Data` methods

func (t *task) Children() []tree.Node {
	var treeNodes []tree.Node
	for _, child := range t.children {
		treeNodes = append(treeNodes, tree.Node(child))
	}
	return treeNodes
}

func (t *task) Data() interface{} {
	return t.String()
}

func printTaskTree(t *task) {
	displayGraph(tree.SprintHrn(t))
}

// PrintTasksTree prints the process tree for all tasks which have
// the specified `field` equals to `value`.
func PrintTasksTree(field, value string) {
	setRenderMode(treeMode)
	selectedTasks := getSelectedTasksFromField(field, value)
	displayGraph(imageTree, fmt.Sprintf("Task Tree for %s: %s", getFullFieldName(field), value))
	for _, task := range selectedTasks {
		printTaskTree(task)
		// tree.SprintHrn generates a white line so we don't need to add an extra one.
	}
}
