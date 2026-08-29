package main

func ValidateRelatedTasks(tasks []Task, relatedTasks []string) error {
	// TODO: implement
	// subtract sets of relatedTasks and tasks.IDs, if any remain, return error
	NotImplemented()
	return nil
}

func SortTasksByIndex(tasks map[string]Task, index TaskIndex) []Task {
	sortedTasks := make([]Task, len(tasks))

	for i := range index.Entries {
		for taskID := range tasks {
			if index.Entries[i].TaskID == taskID {
				sortedTasks = append(sortedTasks, tasks[taskID])
			}
		}
	}

	return sortedTasks
}
