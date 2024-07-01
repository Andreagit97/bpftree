package task

var (
	tasksList []*Task
)

func getTaskFromTid(tid int) *Task {
	for _, task := range tasksList {
		if task.GetTid() == tid {
			return task
		}
	}
	return nil
}

func getSelectedTasksFromField(fieldType allowedField, value string) []*Task {
	var selectedTasks []*Task
	for _, task := range tasksList {
		if taskFieldString(task, fieldType) == value {
			selectedTasks = append(selectedTasks, task)
		}
	}
	return selectedTasks
}
