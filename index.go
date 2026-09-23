package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type TaskIndexEntry struct {
	TaskID   string `json:"id"`
	TaskPath string `json:"path"`
}

type TaskIndex struct {
	// TODO: define interface
	CurrentTaskID string           `json:"current_task"` // could be undefined
	Entries       []TaskIndexEntry `json:"tasks"`
}

func (idx TaskIndex) String() string {
	indexData, err := json.Marshal(
		idx,
		jsontext.WithIndentPrefix(""),
		jsontext.WithIndent("\t"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to serialize task index: %v", err))
	}

	return string(indexData)
}

func LoadIndex(tasksDir string) (*TaskIndex, error) {
	indexFilePath := getIndexFilePath(tasksDir)

	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("task index file `%s` does not exist", indexFilePath)
	}

	indexData, err := os.ReadFile(indexFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read `%s`: %v", indexFilePath, err)
	}

	var index TaskIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		return nil, fmt.Errorf("invalid json")
	}

	return &index, nil
}

func (idx *TaskIndex) AddTask(task Task, taskDir string) {
	idx.Entries = append(idx.Entries, TaskIndexEntry{task.ID, taskDir})
}

func (idx *TaskIndex) RemoveTask(task Task) {
	for i, entry := range idx.Entries {
		if entry.TaskID == task.ID {
			idx.Entries = append(idx.Entries[:i], idx.Entries[i+1:]...)
			break
		}
	}
}

func (idx TaskIndex) SaveIndex(tasksDir string) error {
	indexFilePath := getIndexFilePath(tasksDir)

	file, err := os.OpenFile(indexFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open index file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(idx.String()); err != nil {
		return fmt.Errorf("failed to write to index file: %v", err)
	}

	return nil
}

func (idx TaskIndex) Validate(tasks map[string]Task) error {
	indexEntriesSeq := SliceToSeq(idx.Entries)
	indexIDs := NewSetFrom(indexEntriesSeq, func(entry TaskIndexEntry) string { return entry.TaskID })

	tasksSeq := MapKeysToSeq(tasks)
	tasksIDs := NewSetFrom(tasksSeq, func(taskID string) string { return taskID })

	absentTasks := indexIDs.Difference(*tasksIDs)
	if absentTasks.Len() > 0 {
		return fmt.Errorf("task(s) %v present on index, but failed to be discovered", slices.Collect(absentTasks.Items()))
	}

	absentIndexEntries := tasksIDs.Difference(*indexIDs)
	if absentIndexEntries.Len() > 0 {
		return fmt.Errorf("task(s) %v discovered, but absent in index", slices.Collect(absentIndexEntries.Items()))
	}

	return nil
}

func getIndexFilePath(tasksDir string) string {
	return filepath.Join(tasksDir, "index.json")
}
