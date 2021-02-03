package internal

import (
	"github.com/mlevieux/thodo/src/internal/memory/fsmemory"
	"github.com/mlevieux/thodo/src/internal/todo"
)

type Memory interface {
	GetAllTasks() ([]*todo.Task, error)
	SaveTask(task *todo.Task) (int64, error)
	GetTask(id int64) (*todo.Task, error)
	DeleteTask(id int64) error
}

var (
	_ Memory = &fsmemory.FSMemory{}
)
