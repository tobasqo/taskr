package main

func ValidateRelatedTasks(tasks []Task, relatedTasks []string) error {
	// TODO: implement
	// subtract sets of relatedTasks and tasks.IDs, if any remain, return error
	NotImplemented()
	return nil
}

func SortTasksByIndex(tasks map[string]Task, index TaskIndex) []Task {
	sortedTasks := make([]Task, 0, len(tasks))

	for i := range index.Entries {
		taskID := index.Entries[i].TaskID
		if task, exists := tasks[taskID]; exists {
			sortedTasks = append(sortedTasks, task)
		}
	}

	return sortedTasks
}
