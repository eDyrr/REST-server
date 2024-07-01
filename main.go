package main

import (
	"log"
	"net/http"
)

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

	log.Fatal(http.ListenAndServe("localhost:"+os.Geten("SERVERPORT"), mux))

}