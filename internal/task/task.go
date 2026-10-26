package task

import (
	"slices"
	"time"
)

var timingFormat = "2006-01-02 15:04:05"

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TaskFile struct {
	NextID int    `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

func (taskFile *TaskFile) AddTask(description string) int {
	id := taskFile.NextID
	now := time.Now().Format(timingFormat)

	taskFile.Tasks = append(taskFile.Tasks, Task{
		ID:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	taskFile.NextID++

	return id
}

func (taskFile *TaskFile) DeleteTask(id int) bool {
	oldLength := len(taskFile.Tasks)

	taskFile.Tasks = slices.DeleteFunc(taskFile.Tasks, func(task Task) bool {
		return task.ID == id
	})

	return oldLength != len(taskFile.Tasks)
}

func (taskFile *TaskFile) findIndex(id int) int {
	return slices.IndexFunc(taskFile.Tasks, func(task Task) bool {
		return task.ID == id
	})
}

func (taskFile *TaskFile) Mark(id int, status string) bool {
	taskIndex := taskFile.findIndex(id)

	if taskIndex != -1 {
		taskFile.Tasks[taskIndex].Status = status
		taskFile.Tasks[taskIndex].CreatedAt = time.Now().Format(timingFormat)
		return true
	}
	return false
}

func (taskFile *TaskFile) Update(id int, description string) bool {
	taskIndex := taskFile.findIndex(id)

	if taskIndex != -1 {
		taskFile.Tasks[taskIndex].Description = description
		taskFile.Tasks[taskIndex].UpdatedAt = time.Now().Format(timingFormat)
		return true
	}
	return false
}
