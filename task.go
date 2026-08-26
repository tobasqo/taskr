package main

import (
	"fmt"

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

// TODO: consider removing RelatedTasks section to make the body more flexible with its contents
type taskBody struct {
	Title        string
	Description  string
	RelatedTasks string
}

func (tb taskBody) String() string {
	// TODO: handle empty description and related tasks gracefully
	return fmt.Sprintf(`# %s

%s

## Related Tasks
%s`, tb.Title, tb.Description, tb.RelatedTasks)
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
	return fmt.Sprintf(`%s
---

%s`, t.taskFrontmatter.String(), t.taskBody.String())
}

func PrintTask(task Task, tasksDir string) {
	taskDir := fmt.Sprintf("%s/%s", tasksDir, task.ID)
	fmt.Printf("# %s\npath: %s\n%s", task.Title, taskDir, task.taskFrontmatter.String())
}
