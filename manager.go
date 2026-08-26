package main

import (
	"fmt"
)

type TaskManager interface {
	GetTaskById(taskId string) (*Task, error)
	AddTask(id, title string) (string, error)
}

type LfsTaskManager struct {
	lfsTaskFileManager TaskFileManager
	tasks              []Task // consider using map when many tasks
}

func NewLfsTaskManager(tasksDir string) (*LfsTaskManager, error) {
	taskFileManager := LfsTaskFileManager{
		tasksDir: tasksDir,
	}
	discoveredTasks, err := taskFileManager.DiscoverTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to discover tasks: %v", err)
	}

	tasks := &LfsTaskManager{
		lfsTaskFileManager: taskFileManager,
		tasks:              discoveredTasks,
	}
	return tasks, nil
}

func (tm *LfsTaskManager) GetTaskById(taskId string) (*Task, error) {
	for i := range tm.tasks {
		if tm.tasks[i].ID == taskId {
			return &tm.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID `%s` not found", taskId)
}

func (tm *LfsTaskManager) AddTask(id, title string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("task title is required")
	}

	if !taskIdIsUnique(*tm, id) {
		return "", fmt.Errorf("task `%s` already exists at `%s`", id, tm.lfsTaskFileManager.Location())
	}

	task := NewTask(id, title)
	taskFilePath, err := tm.lfsTaskFileManager.Save(*task)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	tm.tasks = append(tm.tasks, *task)

	return taskFilePath, nil
}

func taskIdIsUnique(tm LfsTaskManager, taskId string) bool {
	return !taskExists(tm, taskId)
}

func taskExists(tm LfsTaskManager, taskId string) bool {
	for i := range tm.tasks {
		if tm.tasks[i].ID == taskId {
			return true
		}
	}
	return false
}
