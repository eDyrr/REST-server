package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type taskServer struct {
	store *TaskStore
	sync.Mutex
}

func NewTaskServer() *taskServer {
	store := New()
	return &taskServer{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	task, err := ts.store.GetTask(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(mux.Vars(req)["year"])
	month, _ := strconv.Atoi(mux.Vars(req)["month"])
	day, _ := strconv.Atoi(mux.Vars(req)["day"])

	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)

	renderJSON(w, tasks)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	tag := mux.Vars(req)["tag"]

	if len(tag) == 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tasks := ts.store.GetTaskByTag(tag)

	renderJSON(w, tasks)
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewTaskServer()
	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.getAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/", server.deleteAllTasksHandler).Methods("DELETE")
	router.HandleFunc("/task/{id:[0-9]+}/", server.getTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}/", server.deleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tag/{tag}/", server.tagHandler).Methods("GET")
	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", server.dueHandler).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}
