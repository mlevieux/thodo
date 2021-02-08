package sqlmemory

import (
	"database/sql"
	"fmt"
	"github.com/mlevieux/thodo/back/internal"
	"github.com/mlevieux/thodo/back/internal/todo"
	"time"
)

type SQLMemory struct {
	db *sql.DB
}

var (
	_ internal.Memory = &SQLMemory{}
)

func dsn(user string, pass string, port string, host string, base string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowCleartextPasswords=1", user, pass, host, port, base)
}

func NewSQLMemory(user string, pass string, port string, host string, base string) (*SQLMemory, error) {

	db, err := sql.Open("mysql", dsn(user, pass, port, host, base))
	if err != nil {
		return nil, err
	}

	sqlMem := &SQLMemory{db: db}
	return sqlMem, nil
}

func (s SQLMemory) GetAllTasks() ([]*todo.Task, error) {

	query := "SELECT id, name, description, state, priority, value, created_at FROM tasks"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*todo.Task, 0)
	for rows.Next() {
		var (
			creationTime string
		)

		t := new(todo.Task)
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.State, &t.Priority, &t.Value, creationTime)
		if err != nil {
			continue
		}

		t.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTime)
		if err != nil {
			continue
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (s SQLMemory) SaveTask(task *todo.Task) (int64, error) {
	query := "INSERT INTO tasks (name, description, state, priority, value) VALUES (?, ?, ?, ?, ?)"
	res, err := s.db.Exec(query, task.Name, task.Description, task.State, task.Priority, task.Value)
	if err != nil {
		return 0, err
	}

	task.Id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return task.Id, err
}

func (s SQLMemory) GetTask(id int64) (*todo.Task, error) {
	query := "SELECT name, description, state, priority, value, created_at FROM tasks WHERE id = ?"

	t := new(todo.Task)
	var creationTime string

	err := s.db.
		QueryRow(query, id).
		Scan(&t.Name, &t.Description, &t.State, &t.Priority, &t.Value, creationTime)
	if err != nil {
		return nil, err
	}

	t.Id = id
	t.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTime)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s SQLMemory) DeleteTask(id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := s.db.Exec(query, id)

	return err
}