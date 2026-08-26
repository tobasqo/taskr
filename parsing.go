package main

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"
)

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
