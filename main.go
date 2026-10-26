package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-cli/internal/storage"
	"task-cli/internal/task"
)

var argLength = len(os.Args)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fatal error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printHelp()
		return nil
	}

	switch os.Args[1] {
	case "-h", "--help":
		printHelp()
		return nil

	case "add":
		err := validateExactArgCount(3, "description not specified", "too many arguments")
		if err != nil {
			return err
		}

		id, err := handleAddTask()
		if err != nil {
			return err
		}
		fmt.Printf("Added task %s: ID = %d\n", os.Args[2], id)
		return nil

	case "update":
		err := validateExactArgCount(4, "id or description not specified", "too many arguments")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
		success, err := handleUpdateTask(id, os.Args[3])
		if err != nil {
			return err
		}
		if !success {
			return fmt.Errorf("task with ID = %d does not exist", id)
		}
		fmt.Printf("Updated task with ID = %d\n", id)
		return nil

	case "delete":
		return handleIDCommand("Deleted task", handleDeleteTask)

	case "mark-in-progress":
		return handleIDCommand(
			"Marked task as in-progress",
			func(id int) (bool, error) {
				return handleMarkTask(id, "in-progress")
			},
		)

	case "mark-done":
		return handleIDCommand(
			"Marked task as done",
			func(id int) (bool, error) {
				return handleMarkTask(id, "done")
			},
		)

	case "list":
		return handleListTasks()

	default:
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func printHelp() {
	lines := []string{
		"",
		"Task tracker CLI",
		"",
		"Usage: task-cli <command> <arguments>",
		"",
		"Available commands:",
		"  add <description> - Add a new task with provided description",
		"  update <id> <description> - Update an existing task by ID with provided description",
		"  delete <id> - Remove an existing task by ID",
		"  mark-in-progress <id> - Mark task as in progress by ID",
		"  mark-done <id> - Mark task as done by ID",
		"  list [done/todo/in-progress] - List tasks by status. Lists all if no status specified",
		"",
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func validateExactArgCount(count int, messageOnLessThan, messageOnMoreThan string) error {
	if argLength < count {
		return errors.New(messageOnLessThan)
	}
	if argLength > count {
		return errors.New(messageOnMoreThan)
	}
	return nil
}

func handleAddTask() (id int, err error) {
	taskFile, err := storage.ReadTasks()
	if err != nil {
		return -1, err
	}
	id = taskFile.AddTask(os.Args[2])

	if err = storage.WriteTasks(taskFile); err != nil {
		return -1, err
	}
	return id, nil
}

func handleIDCommand(successMessage string, action func(int) (bool, error)) error {
	if err := validateExactArgCount(3, "id not specified", "too many arguments"); err != nil {
		return err
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	success, err := action(id)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("task with ID = %d does not exist", id)
	}

	fmt.Printf("%s with ID = %d\n", successMessage, id)
	return nil
}

func handleExistingTask(id int, action func(*task.TaskFile) bool) (bool, error) {
	taskFile, err := storage.ReadTasks()
	if err != nil {
		return false, err
	}

	if id <= 0 || id >= taskFile.NextID {
		return false, fmt.Errorf("task ID out of range")
	}

	success := action(&taskFile)
	if !success {
		return false, nil
	}

	if err := storage.WriteTasks(taskFile); err != nil {
		return false, err
	}

	return true, nil
}

func handleDeleteTask(id int) (bool, error) {
	return handleExistingTask(id, func(taskFile *task.TaskFile) bool {
		return taskFile.DeleteTask(id)
	})
}

func handleMarkTask(id int, status string) (bool, error) {
	return handleExistingTask(id, func(taskFile *task.TaskFile) bool {
		return taskFile.Mark(id, status)
	})
}

func handleUpdateTask(id int, description string) (success bool, err error) {
	return handleExistingTask(id, func(taskFile *task.TaskFile) bool {
		return taskFile.Update(id, description)
	})
}

func fetchListFilter() (filter string, err error) {
	if argLength == 3 {
		filter = os.Args[2]
		switch filter {
		case "todo", "in-progress", "done":
			return filter, nil
		default:
			return "", fmt.Errorf("invalid filter: %s", filter)
		}
	}
	return "", nil
}

func handleListTasks() (err error) {
	if argLength > 3 {
		return fmt.Errorf("too many arguments")
	}

	filter, err := fetchListFilter()
	if err != nil {
		return err
	}

	taskFile, err := storage.ReadTasks()
	if err != nil {
		return err
	}

	filteredTasks := taskFile.Tasks

	if filter != "" {
		filteredTasks = make([]task.Task, 0)
		for _, currentTask := range taskFile.Tasks {
			if currentTask.Status == filter {
				filteredTasks = append(filteredTasks, currentTask)
			}
		}
	}

	for _, currentTask := range filteredTasks {
		fmt.Printf(
			"\nTask %s:\n"+
				"  - id: %d\n"+
				"  - status: %s\n"+
				"  - created at: %s\n"+
				"  - updated at: %s\n\n",
			currentTask.Description,
			currentTask.ID, currentTask.Status,
			currentTask.CreatedAt,
			currentTask.UpdatedAt)
	}
	return nil
}
