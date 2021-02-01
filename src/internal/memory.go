package internal

type Memory interface {
	GetAllTasks() ([]*Task, error)
	SaveTask(task *Task) (int64, error)
	GetTask(id int64) (*Task, error)
	DeleteTask(id int64) error
}

var (
	_ Memory = &FSMemory{}
)
