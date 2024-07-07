package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

// func (ts *taskServer) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	// types used internally in this handler to (de-)serialize the request and
// 	// response from/to JSON
// 	type RequestTask struct {
// 		Text string    `json:"text"`
// 		Tags []string  `json:"tags"`
// 		Due  time.Time `json:"due"`
// 	}

// 	type ResponseId struct {
// 		Id int `json:"id"`
// 	}

// 	// enforce a JSON Content-Type
// 	contentType := req.Header.Get("Content-Type")
// 	mediatype, _, err := mime.ParseMediaType(contentType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if mediatype != "application/json" {
// 		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	dec := json.NewDecoder(req.Body)
// 	dec.DisallowUnknownFields()
// 	var rt RequestTask
// 	if err := dec.Decode(&rt); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	ts.Lock()
// 	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
// 	ts.Unlock()
// 	renderJSON(w, ResponseId{Id: id})
// }

func (ts *taskServer) CreateTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	// assigning `rt` a "task" that came from the frontend via json
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)

	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func (ts *taskServer) deleteAllTasksHandler(c *gin.Context) {
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) deleteTaskHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))

	err := ts.store.DeleteTask(id)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
}

// func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
// 	ts.store.DeleteAllTasks()
// }

// func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	id, _ := strconv.Atoi(mux.Vars(req)["id"])

// 	ts.Lock()
// 	err := ts.store.DeleteTask(id)
// 	ts.Unlock()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// }

// func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling get task at %s\n", req.URL.Path)

// 	id, _ := strconv.Atoi(mux.Vars(req)["id"])

// 	task, err := ts.store.GetTask(id)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	renderJSON(w, task)
// }

func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := ts.store.GetTask(id)

	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

// func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
// 	year, _ := strconv.Atoi(mux.Vars(req)["year"])
// 	month, _ := strconv.Atoi(mux.Vars(req)["month"])
// 	day, _ := strconv.Atoi(mux.Vars(req)["day"])

// 	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)

// 	renderJSON(w, tasks)
// }

func (ts *taskServer) dueHandler(c *gin.Context) {
	year, err := strconv.Atoi(c.Params.ByName("year"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	month, err := strconv.Atoi(c.Params.ByName("month"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	day, err := strconv.Atoi(c.Params.ByName("day"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)

	c.JSON(http.StatusOK, tasks)
}

// func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
// 	tag := mux.Vars(req)["tag"]

// 	if len(tag) == 0 {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	tasks := ts.store.GetTaskByTag(tag)

// 	renderJSON(w, tasks)
// }

func (ts *taskServer) tagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")

	tasks := ts.store.GetTaskByTag(tag)

	c.JSON(http.StatusOK, tasks)
}

// func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
// 	tasks := ts.store.getAllTasks()

// 	if len(tasks) == 0 {
// 		http.Error(w, "no tasks available", http.StatusNotFound)
// 		return
// 	}
// 	renderJSON(w, tasks)
// }

func (ts *taskServer) getAllTasksHandler(c *gin.Context) {
	allTasks := ts.store.getAllTasks()
	c.JSON(http.StatusOK, allTasks)
}

func main() {
	// router := mux.NewRouter()
	// router.StrictSlash(true)
	// server := NewTaskServer()

	// router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	// router.HandleFunc("/task/", server.getAllTasksHandler).Methods("GET")
	// router.HandleFunc("/task/", server.deleteAllTasksHandler).Methods("DELETE")
	// router.HandleFunc("/task/{id:[0-9]+}/", server.getTaskHandler).Methods("GET")
	// router.HandleFunc("/task/{id:[0-9]+}/", server.deleteTaskHandler).Methods("DELETE")
	// router.HandleFunc("/tag/{tag}/", server.tagHandler).Methods("GET")
	// router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", server.dueHandler).Methods("GET")

	// log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), router))
	router := gin.Default()
	server := NewTaskServer()

	router.POST("/task/", server.CreateTaskHandler)
	router.GET("/task/", server.getAllTasksHandler)
	router.DELETE("/task/", server.deleteAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.DELETE("/task/:id", server.deleteTaskHandler)
	router.GET("/tag/:tag", server.tagHandler)
	router.GET("/due/:year/:month/:day", server.dueHandler)
}
