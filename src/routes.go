package main

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"github.com/mlevieux/thodo/src/internal"
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
			task = new(internal.Task)
		)

		err := jsoniter.NewDecoder(r.Body).Decode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task.Apply(internal.WithState(internal.StateTodo))
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