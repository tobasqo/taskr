package main

import "fmt"

type TaskManager struct {
	TasksDir string
	Tasks    *[]Task // consider using map when many tasks
}

func NewTaskManager(tasksDir string) (*TaskManager, error) {
	discoveredTasks, err := DiscoverTasks(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to discover tasks: %v", err)
	}

	tasks := &TaskManager{
		TasksDir: tasksDir,
		Tasks:    discoveredTasks,
	}
	return tasks, nil
}

func (tm *TaskManager) GetTaskById(taskId string) (*Task, error) {
	for _, task := range *tm.Tasks {
		if task.ID == taskId {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task with ID `%s` not found", taskId)
}

func (tm *TaskManager) AddTask(id, title string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("task title is required")
	}

	if !TaskIdIsUnique(id, tm) {
		return "", fmt.Errorf("task `%s` already exists at `%s`", id, tm.TasksDir)
	}

	task := NewTask(id, title)
	taskDir, err := task.Save(tm.TasksDir)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	*tm.Tasks = append(*tm.Tasks, *task)

	return taskDir, nil
}

func (tm *TaskManager) TaskExists(taskId string) bool {
	for _, task := range *tm.Tasks {
		if task.ID == taskId {
			return true
		}
	}
	return false
}

func (tm *TaskManager) CheckValidRelatedTasks(relatedTasks *[]string) error {
	// TODO: implement
	// subtract sets of relatedTasks and tasks.IDs, if any remain, return error
	panic("`checkValidRelatedTasks` not implemented")
}
