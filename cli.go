package main

import (
	"fmt"
	"os"
)

const defaultTasksDir = "./.tasks"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	// TODO: add option to specify tasks directory as optional flag argument
	tasks, err := NewTaskManager(defaultTasksDir)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	switch command {
	case "add":
		taskDir, err := addTask(os.Args[2:], tasks)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		fmt.Printf("Task added at: `%s`\n", taskDir)
	case "list":
		listTasks(os.Args[2:], tasks)
	case "show":
		showTask(os.Args[2:], tasks)
	case "help":
		printHelp()
	default:
		printHelp()
	}
}

func addTask(args []string, tasks *TaskManager) (string, error) {
	// TODO: support specifying task ID and other metadata as named arguments
	if len(args) < 2 {
		return "", fmt.Errorf("task ID and title are required")
	}

	id := args[0]
	title := args[1]
	return tasks.AddTask(id, title)
}

func listTasks(args []string, tasks *TaskManager) {
	// TODO: support filtering tasks by status, related tasks, tags
	// TODO: implement
	panic("`listTasks` not implemented")
}

func showTask(args []string, tasks *TaskManager) {
	id := args[0]
	task, err := tasks.GetTaskById(id)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	task.Display()
}

func printHelp() {
	// TODO: add details for each command and its options
	println(`
Usage:
  taskr <command> [options]

Commands:
  add    Add a new task
  list   List all tasks
  show   Show a specific task
`)
}
