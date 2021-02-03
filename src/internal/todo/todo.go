package todo

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

type TaskState int

const (
	StateTodo TaskState = iota
	StateDoing
	StateCheck
	StateDone
)

type TaskPriority int

const (
	PriorityVeryLow TaskPriority = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityVeryHigh
	PriorityAbsolute
)

type TaskValue int

const (
	ValueNotEstimated TaskValue = iota
	ValueSuperficial
	ValueNeeded
	ValueBlocker
	ValueFinalize
	ValueLarge
	ValueProbablyNecessary
	ValueAbsolute
)

type Task struct {
	Id int64 `json:"Id"`

	Name        string       `json:"name"`
	Description string       `json:"description"`
	State       TaskState    `json:"state"`
	Priority    TaskPriority `json:"priority"`
	Value       TaskValue    `json:"value"`

	CreationTime time.Time `json:"creationTime"`
}

func (task *Task) UnmarshalJSON(data []byte) error {
	type UnmarshallableTask Task

	ut := UnmarshallableTask{}
	err := jsoniter.Unmarshal(data, &ut)
	if err != nil {
		return err
	}

	*task = Task(ut)
	if task.Id == 0 {
		task.Id = -1
	}

	if task.CreationTime.IsZero() {
		task.CreationTime = time.Now()
	}

	return nil
}

type TaskOption func(task *Task)

func WithDescription(desc string) TaskOption {
	return func(task *Task) {
		task.Description = desc
	}
}

func WithState(s TaskState) TaskOption {
	return func(task *Task) {
		task.State = s
	}
}

func WithPriority(prio TaskPriority) TaskOption {
	return func(task *Task) {
		task.Priority = prio
	}
}

func WithValue(val TaskValue) TaskOption {
	return func(task *Task) {
		task.Value = val
	}
}

func NewTask(name string, opts ...TaskOption) *Task {
	t := &Task{
		Id:           -1,
		Name:         name,
		CreationTime: time.Now(),
	}

	t.Apply(opts...)
	return t
}

func (task *Task) Apply(opts ...TaskOption) {
	for _, opt := range opts {
		opt(task)
	}
}

func NewTaskFromPayload(id int64, payload []byte) (*Task, error) {
	task := &Task{
		Id: id,
	}
	err := jsoniter.Unmarshal(payload, task)

	return task, err
}
