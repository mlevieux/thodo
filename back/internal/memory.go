package internal

import (
	"github.com/mlevieux/thodo/back/internal/todo"
)

type TaskStore interface {
	GetAllTasks() ([]*todo.Task, error)
	SaveTask(task *todo.Task) (int64, error)
	GetTask(id int64) (*todo.Task, error)
	DeleteTask(id int64) error
}

type Memory interface {
	TaskStore
}