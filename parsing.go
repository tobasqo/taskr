package main

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"
)

func ParseTaskFromString(taskContent string) (*Task, error) {
	const delimiter = "---"

	if !strings.HasPrefix(taskContent, delimiter) {
		return nil, fmt.Errorf("invalid task format: missing frontmatter or body")
	}

	parts := strings.SplitN(strings.TrimPrefix(taskContent, delimiter), delimiter, 2)
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

func parseTaskFrontmatter(frontmatter string) (*taskFrontmatter, error) {
	frontmatter = strings.Trim(frontmatter, "\r\n ")

	if frontmatter == "" {
		return nil, fmt.Errorf("frontmatter is empty")
	}

	var tf taskFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &tf); err != nil {
		return nil, fmt.Errorf("invalid YAML: %v", err)
	}

	if strings.TrimSpace(tf.ID) == "" {
		return nil, fmt.Errorf("frontmatter field `id` is required")
	}

	if strings.TrimSpace(tf.Status) == "" {
		return nil, fmt.Errorf("frontmatter field `status` is required")
	}

	return &tf, nil
}

func parseTaskBody(body string) (*taskBody, error) {
	body = strings.TrimSpace(body)

	titleLines := strings.SplitN(body, "\n", 2)
	if len(titleLines) < 1 {
		return nil, fmt.Errorf("invalid task body format: missing title")
	}

	var tb taskBody
	tb.Title = strings.TrimPrefix(titleLines[0], "# ")
	if len(titleLines) > 1 {
		tb.Description = strings.TrimSpace(titleLines[1])
	}

	return &tb, nil
}
