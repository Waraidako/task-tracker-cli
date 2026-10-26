package storage

import (
	"encoding/json/v2"
	"errors"
	"os"
	"task-cli/internal/task"
)

func ReadTasks() (taskFile task.TaskFile, err error) {
	taskFile = task.TaskFile{
		NextID: 1,
		Tasks:  make([]task.Task, 0),
	}
	f, err := os.Open("task-list.json")

	if errors.Is(err, os.ErrNotExist) {
		return taskFile, nil
	}

	if err != nil {
		return task.TaskFile{}, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return task.TaskFile{}, err
	}
	if fileInfo.Size() == 0 {
		return taskFile, nil
	}

	if err = json.UnmarshalRead(f, &taskFile); err != nil {
		return task.TaskFile{}, err
	}

	return taskFile, nil
}

func WriteTasks(taskFile task.TaskFile) error {
	tmp, err := os.CreateTemp(".", "task-list-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if err = json.MarshalWrite(tmp, taskFile); err != nil {
		return err
	}

	if err = tmp.Chmod(0o644); err != nil {
		return err
	}

	if err = tmp.Sync(); err != nil {
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), "task-list.json")
}
