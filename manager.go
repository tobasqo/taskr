package main

import (
	"fmt"
)

type TaskManager interface {
	GetTaskById(taskId string) (Task, error)
	AddTask(id, title string) (string, error)
	Tasks() []Task
}

type LfsTaskManager struct {
	lfsTaskFileManager TaskFileManager
	tasks              map[string]Task
	index              TaskIndex
}

func NewLfsTaskManager(tasksDir string) (*LfsTaskManager, error) {
	taskFileManager := LfsTaskFileManager{
		tasksDir: tasksDir,
	}
	discoveredTasks, err := taskFileManager.DiscoverTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to discover tasks: %v", err)
	}

	tasks := make(map[string]Task, len(discoveredTasks))
	for i := range discoveredTasks {
		tasks[discoveredTasks[i].ID] = discoveredTasks[i]
	}

	index, err := LoadIndex(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("could not load index: %v", err)
	}

	if err := index.Validate(tasks); err != nil {
		return nil, fmt.Errorf("error validating index: %v", err)
	}

	taskManager := &LfsTaskManager{
		lfsTaskFileManager: taskFileManager,
		tasks:              tasks,
		index:              *index,
	}
	return taskManager, nil
}

func (tm LfsTaskManager) GetTaskById(taskID string) (Task, error) {
	// probably somehow return *Task if possible?
	// or just add `UpdateTask` method when needed
	task, exists := tm.tasks[taskID]
	if !exists {
		return Task{}, fmt.Errorf("task with ID `%s` not found", taskID)
	}

	return task, nil
}

func (tm *LfsTaskManager) AddTask(id, title string) (string, error) {
	if !taskIdIsUnique(*tm, id) {
		return "", fmt.Errorf("task `%s` already exists at `%s`", id, tm.lfsTaskFileManager.Location())
	}

	task := NewTask(id, title)
	taskFilePath, err := tm.lfsTaskFileManager.Save(*task)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	tm.tasks[id] = *task

	tm.index.AddTask(*task, tm.lfsTaskFileManager.GetTaskDir(id))
	tm.index.Validate(tm.tasks)
	// not optimal to save index every time we add a task
	tm.index.SaveIndex(tm.lfsTaskFileManager.Location())

	return taskFilePath, nil
}

func (tm *LfsTaskManager) Tasks() []Task {
	return SortTasksByIndex(tm.tasks, tm.index)
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
