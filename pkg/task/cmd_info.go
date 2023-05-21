package task

import "fmt"

func printTaskInfo(t *task) {
	displayGraph(imageNewspaper, t)
}

// PrintTasksInfo prints info about tasks which have the specified `field` equals to `value`.
func PrintTasksInfo(field, value string) {
	setRenderMode(infoMode)
	selectedTasks := getSelectedTasksFromField(field, value)
	displayGraph(imageInfo, fmt.Sprintf("Task Info for %s: %s", getFullFieldName(field), value))
	for _, task := range selectedTasks {
		printTaskInfo(task)
	}
}
