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

type TaskStore map[int]Task

func New() *TaskStore {
	taskstore := &TaskStore{}
	return taskstore
}

func (ts TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	newTask := Task{Id: len(ts), Text: text, Tags: tags, Due: due}
	ts[newTask.Id] = newTask
	return newTask.Id
}

func (ts TaskStore) GetTask(id int) (Task, error) {
	task, found := ts[id]
	if found {
		return task, nil
	}
	return Task{}, errors.New("task not found")
}

func (ts TaskStore) DeleteTask(id int) error {
	if _, err := ts.GetTask(id); err != nil {
		return nil
	}
	delete(ts, id)
	return nil
}

func (ts TaskStore) DeleteAllTasks() error {
	if len(ts) == 0 {
		return nil
	}
	for _, task := range ts {
		delete(ts, task.Id)
	}
	return nil
}

func (ts TaskStore) getAllTasks() []Task {
	var tasks []Task
	for index, _ := range ts {
		tasks[index] = ts[index]
	}
	return tasks
}

func (ts TaskStore) GetTaskByTag(tag string) []Task {
	var tasks []Task
	for _, task := range ts {
		for _, t := range task.Tags {
			if t == tag {
				tasks = append(tasks, task)
				break
			}
		}
	}
	return tasks
}

func (ts TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task {
	var tasks []Task
	for _, task := range ts {
		if task.Due.Year() == year && task.Due.Month() == month && task.Due.Day() == day {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
