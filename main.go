package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type taskServer struct {
	store *TaskStore
}

func NewTaskServer() *taskServer {
	store := New()
	return &taskServer{store: store}
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req http.Request) {
	year, err := strconv.Atoi(req.PathValue("year"))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(req.PathValue("month"))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	day, err := strconv.Atoi(req.PathValue("day"))
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)

	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	tag := req.PathValue("tag")
	if len(tag) == 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tasks := ts.store.GetTaskByTag(tag)

	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("POST /task/", server.CreateTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	mux.HandleFunc("DELETE /task/", server.deleteAllTasksHandler)
	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	mux.HandleFunc("DELETE /task/{id}/", server.deleteTaskHandler)
	mux.HandleFunc("GET /tag/{tag}/", server.tagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}
