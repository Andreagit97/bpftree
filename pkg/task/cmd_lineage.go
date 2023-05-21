package task

import "fmt"

func printTaskLineage(t *task) {
	for {
		/* Print current task info */
		displayGraph(imageDownArrow, t)

		/* We reached the last ancestor, stop */
		if t.getRealParentTid() == 0 {
			break
		}

		/* go to the real parent */
		t = getTaskFromTid(t.getRealParentTid())
	}
}

// PrintTasksLineage prints lineage of tasks which have the specified `field` equals to `value`.
func PrintTasksLineage(field, value string) {
	setRenderMode(lineageMode)
	selectedTasks := getSelectedTasksFromField(field, value)
	displayGraph(imageLineage, fmt.Sprintf("Task Lineage for %s: %s", getFullFieldName(field), value))
	for _, task := range selectedTasks {
		printTaskLineage(task)
		displayGraph() // Leave a white line between a lineage and another
	}
}
