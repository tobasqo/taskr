package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type taskFrontmatter struct {
	ID           string   `yaml:"id"`
	Status       string   `yaml:"status"`
	RelatedTasks []string `yaml:"related_tasks"` // can be undefined
	Tags         []string `yaml:"tags"`          // can be undefined
}

func (tf *taskFrontmatter) String() string {
	data, err := yaml.Dump(tf, yaml.WithIndent(2))
	if err != nil {
		panic(fmt.Sprintf("failed to serialize task frontmatter: %v", err))
	}
	return string(data)
}

type taskBody struct {
	Title        string
	Description  string
	RelatedTasks string
}

func (tb *taskBody) String() string {
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
		taskFrontmatter: taskFrontmatter{
			ID:     id,
			Status: "todo",
		},
		taskBody: taskBody{
			Title: title,
		},
	}
}

func (t *Task) String() string {
	return fmt.Sprintf(`%s
---

%s`, t.taskFrontmatter.String(), t.taskBody.String())
}

func (t *Task) Display() {
	fmt.Printf("# %s\n%s", t.Title, t.taskFrontmatter.String())
}

func (t *Task) Save(tasksDir string) (string, error) {
	taskDir := fmt.Sprintf("%s/%s", tasksDir, t.ID)

	if _, err := os.Stat(taskDir); os.IsNotExist(err) {
		err := CreateTaskFile(t, taskDir)
		if err != nil {
			return "", err
		}
	}

	return taskDir, nil
}
