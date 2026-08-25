package main

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

func LoadTask(taskPath string) (*Task, error) {
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

func DiscoverTasks(tasksDir string) (*[]Task, error) {
	if _, err := os.Stat(tasksDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("tasks directory `%s` does not exist", tasksDir)
	}
	
	dirEntries, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %v", err)
	}
	
	var tasks []Task
	for _, entry := range dirEntries {
		if entry.IsDir() {
			taskPath := fmt.Sprintf("%s/%s/TASK.md", tasksDir, entry.Name())
			task, err := LoadTask(taskPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load task from `%s`: %v", taskPath, err)
			}
			tasks = append(tasks, *task)
		}
	}

	return &tasks, nil
}

func ParseTaskFromString(taskContent string) (*Task, error) {
	parts := strings.SplitN(taskContent, "---", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid task format: missing frontmatter or body")
	}

	frontmatter := parts[0]
	body := parts[1]

	taskFrontmatter, err := parseTaskFrontmatter(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task frontmatter: %v", err)
	}

	taskBody, err := parseTaskBody(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task body: %v", err)
	}

	task := &Task{
		taskFrontmatter: *taskFrontmatter,
		taskBody:        *taskBody,
	}
	return task, nil
}

func TaskIdIsUnique(taskId string, taskManager *TaskManager) bool {
	return !taskManager.TaskExists(taskId)
}

// TODO: think about renaming this function
func CreateTaskFile(task *Task, taskDir string) error {
	taskFilePath := fmt.Sprintf("%s/TASK.md", taskDir)

	err := os.MkdirAll(taskDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create task directory: %v", err)
	}

	file, err := os.Create(taskFilePath)
	if err != nil {
		return fmt.Errorf("failed to create task file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(task.String())
	if err != nil {
		return fmt.Errorf("failed to write to task file: %v", err)
	}

	return nil
}

func parseTaskFrontmatter(frontmatter string) (*taskFrontmatter, error) {
	// is the trim even needed?
	frontmatter = strings.Trim(frontmatter, "\r\n ")

	var tf taskFrontmatter

	err := yaml.Unmarshal([]byte(frontmatter), &tf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %v", err)
	}

	return &tf, nil
}

func parseTaskBody(body string) (*taskBody, error) {
	body = strings.TrimSpace(body)

	bodyParts := strings.SplitN(body, "## Related Tasks", 2)
	if len(bodyParts) != 2 {
		return nil, fmt.Errorf("invalid task body format: missing '## Related Tasks' section")
	}

	titleAndDescription := strings.TrimSpace(bodyParts[0])
	relatedTasks := strings.TrimSpace(bodyParts[1])

	titleLines := strings.SplitN(titleAndDescription, "\n", 2)
	if len(titleLines) < 1 {
		return nil, fmt.Errorf("invalid task body format: missing title or description")
	}

	var tb taskBody
	tb.Title = strings.TrimPrefix(titleLines[0], "# ")
	if len(titleLines) > 1 {
		tb.Description = strings.TrimSpace(titleLines[1])
	}
	tb.RelatedTasks = relatedTasks

	return &tb, nil
}
