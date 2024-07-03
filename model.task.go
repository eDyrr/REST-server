package main

import (
	"errors"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	store map[int]Task
}

func New() *TaskStore {
	taskstore := &TaskStore{}
	return taskstore
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	newTask := Task{Id: len(ts.store), Text: text, Tags: tags, Due: due}
	ts.store[newTask.Id] = newTask
	return newTask.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	task, found := ts.store[id]
	if found {
		return task, nil
	}
	return Task{}, errors.New("task not found")
}

func (ts *TaskStore) DeleteTask(id int) error {
	if _, err := ts.GetTask(id); err != nil {
		return nil
	}
	delete(ts.store, id)
	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	if len(ts.store) == 0 {
		return nil
	}
	for _, task := range ts.store {
		delete(ts.store, task.Id)
	}
	return nil
}

func (ts *TaskStore) getAllTasks() []Task {
	var tasks []Task
	for index, _ := range ts.store {
		tasks[index] = ts.store[index]
	}
	return tasks
}
