package fsmemory

import (
	"github.com/mlevieux/thodo/src/internal/todo"
	"github.com/sanity-io/litter"
	"os"
	"path/filepath"
	"testing"
)

var (
	testTaskSet []*todo.Task
)

func makeTestTaskSet() {
	testTaskSet = []*todo.Task{
		todo.NewTask(
			"task 1",
			todo.WithDescription("some test task"),
			todo.WithPriority(todo.PriorityNormal),
			todo.WithState(todo.StateTodo),
			todo.WithValue(todo.ValueNotEstimated),
		),
		todo.NewTask(
			"task 2",
			todo.WithDescription("some test task 2"),
			todo.WithPriority(todo.PriorityVeryHigh),
			todo.WithState(todo.StateDoing),
			todo.WithValue(todo.ValueNeeded),
		),
	}
}

func makeRootDir() string {
	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, "taskmem")

	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	return path
}

func TestFSMemory_SaveTask(t *testing.T) {
	mem, err := NewFSMemory(makeRootDir())
	if err != nil {
		panic(err)
	}

	insertTasks(t, mem)
}

func insertTasks(t *testing.T, mem *FSMemory) {
	makeTestTaskSet()
	for _, task := range testTaskSet {
		id, err := mem.SaveTask(task)
		if err != nil {
			panic(err)
		}

		if id != mem.lastTaskId {
			t.Log("LastTaskId =", mem.lastTaskId, "instead of", task.Id)
			t.Fail()
		}
	}
}

func TestFSMemory_GetTask(t *testing.T) {
	mem, err := NewFSMemory(makeRootDir())
	if err != nil {
		panic(err)
	}

	insertTasks(t, mem)
	task, err := mem.GetTask(1)
	if err != nil {
		panic(err)
	}

	if task.Name != "task 1" {
		t.Fail()
	}
}

func TestFSMemory_DeleteTask(t *testing.T) {
	mem, err := NewFSMemory(makeRootDir())
	if err != nil {
		panic(err)
	}

	insertTasks(t, mem)
	err = mem.DeleteTask(1)
	if err != nil {
		t.Fail()
	}
}

func TestFSMemory_GetAllTasks(t *testing.T) {
	mem, err := NewFSMemory(makeRootDir())
	if err != nil {
		panic(err)
	}

	insertTasks(t, mem)
	ts, err := mem.GetAllTasks()
	if err != nil {
		panic(err)
	}

	if len(ts) != len(testTaskSet) {
		t.Log(litter.Sdump(ts))
		t.Fail()
	}
}
