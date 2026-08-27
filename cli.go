package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const defaultTasksDir = "./.tasks"

var commands = []struct {
	name        string
	description string
}{
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
		fmt.Fprintln(os.Stderr, err)

		// TODO: this looks weird
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

	case "list":
		listTasks(args[1:], taskManager)

	case "show":
		if err := showTask(args[1:], taskManager, *tasksDir); err != nil {
			return err
		}
	}

	return nil
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

func listTasks(args []string, taskManager TaskManager) {
	// TODO: support filtering tasks by status, related tasks, tags
	// TODO: implement
	panic("`listTasks` not implemented")
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

	task, err := taskManager.GetTaskById(*id)
	if err != nil {
		return err
	}

	PrintTask(*task, tasksDir)
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
