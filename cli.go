package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

const defaultTasksDir = "./.tasks"

var commands = []struct {
	name        string
	description string
}{
	{name: "init", description: "Initialize tasks directory"},
	{name: "add", description: "Add a new task"},
	{name: "list", description: "List all tasks"},
	{name: "show", description: "Show a specific task"},
	{name: "help", description: "Print help message"},
}

type missingRequiredFlagError struct {
	flagName string
	FlagSet  *flag.FlagSet
}

func (e *missingRequiredFlagError) Error() string {
	return fmt.Sprintf("flag -%s is required", e.flagName)
}

type unexpectedPositionalArgsError struct {
	args    []string
	FlagSet *flag.FlagSet
}

func (e *unexpectedPositionalArgsError) Error() string {
	return fmt.Sprintf("unexpected positional arguments: %v", e.args)
}

func main() {
	if err := run(); err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, "Error:", err)

		// this looks weird
		if missingRequiredFlagErr, ok := errors.AsType[*missingRequiredFlagError](err); ok {
			missingRequiredFlagErr.FlagSet.Usage()
		} else if unexpectedPositionalArgsErr, ok := errors.AsType[*unexpectedPositionalArgsError](err); ok {
			unexpectedPositionalArgsErr.FlagSet.Usage()
		}

		os.Exit(1)
	}
}

func run() error {
	globalFlags := newRootFlagSet()
	tasksDir := globalFlags.String("tasks-dir", defaultTasksDir, "directory containing task files")

	if len(os.Args) < 2 {
		globalFlags.Usage()
		return nil
	}

	if err := globalFlags.Parse(os.Args[1:]); err != nil {
		return err
	}

	args := globalFlags.Args()
	if len(args) == 0 {
		globalFlags.Usage()
		return nil
	}

	command := args[0]
	if command == "init" {
		tasksDir, err := initTasksDirectory(args[1:])
		if err != nil {
			return err
		}
		fmt.Printf("Initialized tasks directory at: `%s`\n", tasksDir)
	} else {
		var taskManager TaskManager
		taskManager, err := NewLfsTaskManager(*tasksDir)
		if err != nil {
			return err
		}

		switch command {
		case "add":
			taskFilePath, err := addTask(args[1:], taskManager)
			if err != nil {
				return err
			}
			fmt.Printf("Task added at: `%s`\n", taskFilePath)

		case "remove":
			taskFilePath, err := removeTask(args[1:], taskManager)
			if err != nil {
				return err
			}
			fmt.Printf("Task at `%s` removed\n", taskFilePath)

		case "list":
			if err := listTasks(args[1:], taskManager, *tasksDir); err != nil {
				return err
			}

		case "show":
			if err := showTask(args[1:], taskManager, *tasksDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func initTasksDirectory(args []string) (string, error) {
	flags := newFlagSet("init")

	dstDir := flags.String("dir", ".", "Destination directory to initialize")

	if err := flags.Parse(args); err != nil {
		return "", err
	}

	if flags.NArg() != 0 && *dstDir != "" {
		*dstDir = args[0]
	}

	dirStat, err := os.Stat(*dstDir)
	if err != nil {
		return "", fmt.Errorf("directory `%s` does not exist", *dstDir)
	}

	if !dirStat.IsDir() {
		return "", fmt.Errorf("provided path `%s` is not a directory", *dstDir)
	}

	tasksDir := fmt.Sprintf("%s/.tasks", *dstDir)
	if _, err := os.Stat(tasksDir); err == nil {
		dirEntries, err := os.ReadDir(tasksDir)
		if err != nil {
			return "", fmt.Errorf("could not read directory `%s`: %v", tasksDir, err)
		}

		if len(dirEntries) > 0 {
			return "", fmt.Errorf("directory `%s` is not empty", tasksDir)
		}
	} else if os.IsNotExist(err) {
		if err := os.Mkdir(tasksDir, 0o755); err != nil {
			return "", fmt.Errorf("could not create directory `%s`: %v", tasksDir, err)
		}
	} else {
		return "", fmt.Errorf("could not stat directory `%s`: %v", tasksDir, err)
	}

	index := TaskIndex{Entries: []TaskIndexEntry{}}
	if err := index.SaveIndex(tasksDir); err != nil {
		return "", fmt.Errorf("failed to save task index: %v", err)
	}

	return tasksDir, nil
}

func addTask(args []string, taskManager TaskManager) (string, error) {
	flags := newFlagSet("add")

	id := flags.String("id", "", "task ID (required)")
	title := flags.String("title", "", "task title (required)")

	if err := flags.Parse(args); err != nil {
		return "", err
	}

	if flags.NArg() != 0 {
		return "", &unexpectedPositionalArgsError{flags.Args(), flags}
	}

	if *id == "" {
		return "", &missingRequiredFlagError{"id", flags}
	}

	if *title == "" {
		return "", &missingRequiredFlagError{"title", flags}
	}

	return taskManager.AddTask(*id, *title)
}

func removeTask(args []string, taskManager TaskManager) (string, error) {
	flags := newFlagSet("remove")

	id := flags.String("id", "", "task ID (required)")

	if err := flags.Parse(args); err != nil {
		return "", err
	}

	if flags.NArg() != 0 {
		return "", &unexpectedPositionalArgsError{flags.Args(), flags}
	}

	if *id == "" {
		return "", &missingRequiredFlagError{"id", flags}
	}

	return taskManager.RemoveTask(*id)
}

func listTasks(args []string, taskManager TaskManager, tasksDir string) error {
	tasks := taskManager.Tasks()

	if len(tasks) == 0 {
		fmt.Printf("No tasks at `%s`\n", tasksDir)
		return nil
	}

	flags := newFlagSet("list")

	status := flags.String("status", "", "task status (optional)")
	// TODO: support filtering tasks by multiple related tasks - should it be OR or AND?
	relatedTask := flags.String("related-task", "", "related task (optional)")
	// TODO: support filtering tasks by multiple tags - should it be OR or AND?
	tag := flags.String("tag", "", "task tag (optional)")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return &unexpectedPositionalArgsError{flags.Args(), flags}
	}

	if *status != "" {
		tasks = filterTasks(tasks, func(task Task) bool { return task.Status == *status })
	}

	if *relatedTask != "" {
		tasks = filterTasks(tasks, func(task Task) bool {
			return slices.Contains(task.RelatedTasks, *relatedTask)
		})
	}

	if *tag != "" {
		tasks = filterTasks(tasks, func(task Task) bool {
			return slices.Contains(task.Tags, *tag)
		})
	}

	if len(tasks) == 0 {
		fmt.Printf("No tasks for given query at `%s`\n", tasksDir)
		return nil
	}

	for i := range tasks {
		println(strings.Repeat("=", 80))
		PrintTask(tasks[i], tasksDir)
	}
	println(strings.Repeat("=", 80))

	return nil
}

func filterTasks(tasks []Task, predicate func(Task) bool) []Task {
	result := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		if predicate(task) {
			result = append(result, task)
		}
	}

	return result
}

func showTask(args []string, taskManager TaskManager, tasksDir string) error {
	flags := newFlagSet("show")

	id := flags.String("id", "", "task ID (required)")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() != 0 {
		return &unexpectedPositionalArgsError{flags.Args(), flags}
	}

	if *id == "" {
		return &missingRequiredFlagError{"id", flags}
	}

	task, err := taskManager.GetTaskByID(*id)
	if err != nil {
		return err
	}

	PrintTask(task, tasksDir)
	return nil
}

func newFlagSet(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ExitOnError)

	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage:\n  taskr %s [options]\n\nOptions:\n", name)
		flags.PrintDefaults()
	}

	return flags
}

func newRootFlagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("taskr", flag.ContinueOnError)

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage:")
		fmt.Fprintln(flags.Output(), "  taskr [--tasks-dir <dir>] <command> [options]")

		fmt.Fprintln(flags.Output(), "\nOptions:")
		flags.PrintDefaults()

		fmt.Fprintln(flags.Output(), "\nCommands:")
		for _, command := range commands {
			fmt.Fprintf(flags.Output(), "  %-6s %s\n", command.name, command.description)
		}
	}
	return flags
}
