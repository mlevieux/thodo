package main

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/mlevieux/thodo/back/internal"
	"github.com/mlevieux/thodo/back/internal/todo"
	"net/http"
)

type handler struct {
	mem internal.Memory
}

func (h handler) getTasks(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:

		tasks, err := h.mem.GetAllTasks()
		if err != nil {
			http.Error(w, "Could not retrieve tasks", http.StatusInternalServerError)
			return
		}

		err = jsoniter.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func (h handler) postTask(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		var (
			task = new(todo.Task)
		)

		err := jsoniter.NewDecoder(r.Body).Decode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task.Apply(todo.WithState(todo.StateTodo))
		id, err := h.mem.SaveTask(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = jsoniter.NewEncoder(w).Encode(struct {
			Id int64 `json:"id"`
		}{
			Id: id,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
