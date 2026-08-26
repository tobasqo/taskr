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

	var tm TaskManager
	tm, err := NewLfsTaskManager(defaultTasksDir)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	switch command {
	case "add":
		taskDir, err := addTask(os.Args[2:], tm)
		if err != nil {
			println("Error:", err.Error())
			return
		}
		fmt.Printf("Task added at: `%s`\n", taskDir)
	case "list":
		listTasks(os.Args[2:], tm)
	case "show":
		showTask(os.Args[2:], tm)
	case "help":
		printHelp()
	default:
		printHelp()
	}
}

func addTask(args []string, taskManager TaskManager) (string, error) {
	// TODO: support specifying task ID and other metadata as named arguments
	if len(args) < 2 {
		return "", fmt.Errorf("task ID and title are required")
	}

	id := args[0]
	title := args[1]
	return taskManager.AddTask(id, title)
}

func listTasks(args []string, taskManager TaskManager) {
	// TODO: support filtering tasks by status, related tasks, tags
	// TODO: implement
	panic("`listTasks` not implemented")
}

func showTask(args []string, taskManager TaskManager) {
	id := args[0]
	task, err := taskManager.GetTaskById(id)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// TODO: don't hardcode defaultTasksDir in here
	PrintTask(*task, defaultTasksDir)
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
