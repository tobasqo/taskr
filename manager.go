package main

import (
	"fmt"
	"path/filepath"
)

type TaskManager interface {
	GetTaskByID(taskID string) (Task, error)
	AddTask(id, title string) (string, error)
	RemoveTask(id string) (string, error)
	Tasks() []Task
}

type LfsTaskManager struct {
	lfsTaskFileManager TaskFileManager
	tasks              map[string]Task
	index              TaskIndex
}

func NewLfsTaskManager(tasksDir string) (*LfsTaskManager, error) {
	tasksDirFullpath, err := filepath.Abs(tasksDir)
	if err != nil {
		return nil, err
	}

	taskFileManager := LfsTaskFileManager{
		tasksDir: tasksDirFullpath,
	}
	discoveredTasks, err := taskFileManager.DiscoverTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to discover tasks: %v", err)
	}

	tasks := make(map[string]Task, len(discoveredTasks))
	for i := range discoveredTasks {
		tasks[discoveredTasks[i].ID] = discoveredTasks[i]
	}

	index, err := LoadIndex(tasksDirFullpath)
	if err != nil {
		return nil, fmt.Errorf("could not load index: %v - forgot to initialize?", err)
	}

	if err := index.Validate(tasks); err != nil {
		return nil, fmt.Errorf("invalid index: %v", err)
	}

	taskManager := &LfsTaskManager{
		lfsTaskFileManager: taskFileManager,
		tasks:              tasks,
		index:              *index,
	}
	return taskManager, nil
}

func (tm LfsTaskManager) GetTaskByID(taskID string) (Task, error) {
	// probably somehow return *Task if possible?
	// or just add `UpdateTask` method when needed
	task, exists := tm.tasks[taskID]
	if !exists {
		return Task{}, fmt.Errorf("task with ID `%s` not found", taskID)
	}

	return task, nil
}

func (tm *LfsTaskManager) AddTask(id, title string) (string, error) {
	if !taskIDIsUnique(*tm, id) {
		return "", fmt.Errorf("task `%s` already exists at `%s`", id, tm.lfsTaskFileManager.Location())
	}

	task := NewTask(id, title)
	taskFilePath, err := tm.lfsTaskFileManager.Save(*task)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	tm.tasks[id] = *task

	tm.index.AddTask(*task, tm.lfsTaskFileManager.GetTaskDir(id))
	if err := tm.index.Validate(tm.tasks); err != nil {
		return "", fmt.Errorf("invalid index: %v", err)
	}
	// not optimal to save index every time we add a task
	if err := tm.index.SaveIndex(tm.lfsTaskFileManager.Location()); err != nil {
		return "", fmt.Errorf("failed to save index: %v", err)
	}

	return taskFilePath, nil
}

func (tm *LfsTaskManager) RemoveTask(id string) (string, error) {
	task, exists := tm.tasks[id]
	if !exists {
		return "", fmt.Errorf("task with ID `%s` not found", id)
	}

	tm.index.RemoveTask(task)
	// not optimal to save index every time we remove a task
	if err := tm.index.SaveIndex(tm.lfsTaskFileManager.Location()); err != nil {
		return "", fmt.Errorf("failed to save index: %v", err)
	}

	taskFilePath, err := tm.lfsTaskFileManager.Delete(task)
	if err != nil {
		return "", fmt.Errorf("failed to delete task: %v", err)
	}

	delete(tm.tasks, id)

	if err := tm.index.Validate(tm.tasks); err != nil {
		return "", fmt.Errorf("invalid index: %v", err)
	}

	return taskFilePath, nil
}

func (tm *LfsTaskManager) Tasks() []Task {
	return SortTasksByIndex(tm.tasks, tm.index)
}

func taskIDIsUnique(tm LfsTaskManager, taskID string) bool {
	return !taskExists(tm, taskID)
}

func taskExists(tm LfsTaskManager, taskID string) bool {
	for i := range tm.tasks {
		if tm.tasks[i].ID == taskID {
			return true
		}
	}
	return false
}
