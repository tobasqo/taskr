package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type TaskFileManager interface {
	GetTaskDir(taskID string) string
	Save(task Task) (string, error)
	Delete(task Task) (string, error)
	LoadTask(taskPath string) (*Task, error)
	DiscoverTasks() ([]Task, error)
	Location() string
}

type LfsTaskFileManager struct {
	tasksDir string
}

func (tfm LfsTaskFileManager) GetTaskDir(taskID string) string {
	return filepath.Join(tfm.tasksDir, taskID)
}

func (tfm LfsTaskFileManager) Save(task Task) (string, error) {
	taskDir := tfm.GetTaskDir(task.ID)

	if err := tfm.ensureTaskDirExists(taskDir); err != nil {
		return "", err
	}

	taskFilePath := filepath.Join(taskDir, "TASK.md")

	file, err := os.OpenFile(taskFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return "", fmt.Errorf("failed to create task file: %v", err)
	}
	defer file.Close()

	if _, err = file.WriteString(task.String()); err != nil {
		return "", fmt.Errorf("failed to write to task file: %v", err)
	}

	return taskFilePath, nil
}

func (tfm LfsTaskFileManager) Delete(task Task) (string, error) {
	taskDir := tfm.GetTaskDir(task.ID)

	if _, err := os.Stat(taskDir); os.IsNotExist(err) {
		return "", fmt.Errorf("task directory `%s` does not exist", taskDir)
	}

	if err := os.RemoveAll(taskDir); err != nil {
		return "", fmt.Errorf("failed to delete task directory: %v", err)
	}

	return taskDir, nil
}

func (tfm LfsTaskFileManager) LoadTask(taskPath string) (*Task, error) {
	if _, err := os.Stat(taskPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("task file `%s` does not exist", taskPath)
	}

	content, err := os.ReadFile(taskPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %v", err)
	}

	task, err := ParseTaskFromString(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse task file: %v", err)
	}

	return task, nil
}

func (tfm LfsTaskFileManager) DiscoverTasks() ([]Task, error) {
	if _, err := os.Stat(tfm.tasksDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("tasks directory `%s` does not exist", tfm.tasksDir)
	}

	dirEntries, err := os.ReadDir(tfm.tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %v", err)
	}

	var tasks []Task
	for _, entry := range dirEntries {
		if entry.IsDir() {
			taskPath := filepath.Join(tfm.tasksDir, entry.Name(), "TASK.md")

			if _, err := os.Stat(taskPath); os.IsNotExist(err) {
				// NOTE: should we print something? probably only if `verbose` or something
				continue
			}

			task, err := tfm.LoadTask(taskPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load task from `%s`: %v", taskPath, err)
			}

			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

func (tfm LfsTaskFileManager) Location() string {
	return tfm.tasksDir
}

func (tfm LfsTaskFileManager) ensureTaskDirExists(taskDir string) error {
	if err := os.MkdirAll(taskDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create task directory: %v", err)
	}

	return nil
}
