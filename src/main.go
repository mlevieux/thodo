package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/mlevieux/thodo/src/internal"
)

func main() {

	r := mux.NewRouter()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mem, err := internal.NewFSMemory(filepath.Join(wd, ".tasks"))
	log.Println("Ho!")
	if err != nil {
		panic(err)
	}

	h := handler{
		mem: mem,
	}

	r.HandleFunc("/tasks", h.getTasks).Methods(http.MethodGet)
	r.HandleFunc("/task", h.postTask).Methods(http.MethodPost)

	log.Fatalln(http.ListenAndServe(":8080", r))
}
