package main

import (
	"fmt"
	"os"
)

type TaskFileManager interface {
	GetTaskDir(taskId string) string
	Save(task Task) (string, error)
	LoadTask(taskPath string) (*Task, error)
	DiscoverTasks() ([]Task, error)
	Location() string
}

type LfsTaskFileManager struct {
	tasksDir string
}

func (tfm LfsTaskFileManager) GetTaskDir(taskId string) string {
	return fmt.Sprintf("%s/%s", tfm.tasksDir, taskId)
}

func (tfm LfsTaskFileManager) Save(task Task) (string, error) {
	taskDir := tfm.GetTaskDir(task.ID)

	err := tfm.ensureTaskDirExists(taskDir)
	if err != nil {
		return "", err
	}

	taskFilePath := fmt.Sprintf("%s/TASK.md", taskDir)

	file, err := os.OpenFile(taskFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create task file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(task.String())
	if err != nil {
		return "", fmt.Errorf("failed to write to task file: %v", err)
	}

	return taskFilePath, nil
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
			taskPath := fmt.Sprintf("%s/%s/TASK.md", tfm.tasksDir, entry.Name())
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
	err := os.MkdirAll(taskDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create task directory: %v", err)
	}

	return nil
}
