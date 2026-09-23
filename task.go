package main

import (
	"fmt"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type taskFrontmatter struct {
	ID           string   `yaml:"id"`
	Status       string   `yaml:"status"`
	RelatedTasks []string `yaml:"related_tasks"` // can be undefined
	Tags         []string `yaml:"tags"`          // can be undefined
}

func (tf taskFrontmatter) String() string {
	data, err := yaml.Dump(tf, yaml.WithIndent(2))
	if err != nil {
		panic(fmt.Sprintf("failed to serialize task frontmatter: %v", err))
	}
	return string(data)
}

type taskBody struct {
	Title       string
	Description string
}

func (tb taskBody) String() string {
	str := fmt.Sprintf("# %s\n", tb.Title)

	if tb.Description == "" {
		return str
	}

	return fmt.Sprintf("%s\n%s", str, tb.Description)
}

type Task struct {
	taskFrontmatter
	taskBody
}

func NewTask(id, title string) *Task {
	return &Task{
		ID:     id,
		Status: "todo",
		Title:  title,
	}
}

func (t Task) String() string {
	return fmt.Sprintf(`---
%s---

%s`, t.taskFrontmatter.String(), t.taskBody.String())
}

func PrintTask(task Task, tasksDir string) {
	taskDir := filepath.Join(tasksDir, task.ID)
	fmt.Printf("# %s\npath: %s\n%s", task.Title, taskDir, task.taskFrontmatter.String())
}
