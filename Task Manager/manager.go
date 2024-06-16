package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	Id     string `json:"TaskID"`
	Status string `json:"TaskStatus"`
}

type taskMap map[string]string

/*Map с задачами*/
var tasksRep = make(taskMap)

/*Запись новой задачи*/
func (t taskMap) post(w http.ResponseWriter, r *http.Request) {
	var taskDecode Task
	err := json.NewDecoder(r.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
	taskDecode.Status = "Proccessing"
	//taskServ = append(taskServ, taskDecode)
	t[taskDecode.Id] = taskDecode.Status
	fmt.Fprintf(w, "Task was added '%s'", taskDecode.Id)
}

/*Получение задачи по ID*/
func (t taskMap) getById(id string, w http.ResponseWriter) {
	value, ok := t[id]
	if ok {
		fmt.Fprintf(w, "Task ID: %s\nStatus: %s", id, value)
	} else {
		fmt.Fprintln(w, "Task not found")
	}

}

/*Получение всех задач*/
func (t taskMap) getAll(w http.ResponseWriter) {
	fmt.Fprintf(w, "%d tasks added \n", len(t))
	if len(t) != 0 {
		json.NewEncoder(w).Encode(t)
	}
}

/*Изменение статуса задачи по ID*/
func (t taskMap) changeStatus(w http.ResponseWriter, r *http.Request) {
	var taskDecode Task
	err := json.NewDecoder(r.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
	_, ok := t[taskDecode.Id]
	if ok {
		t[taskDecode.Id] = taskDecode.Status
		fmt.Fprintf(w, "Status was updated:\nTask ID: %s\nStatus: %s", taskDecode.Id, t[taskDecode.Id])
	} else {
		fmt.Fprintln(w, "Task not found")
	}
}

/*Удаление задачи по ID*/
func (t taskMap) deleteById(id string, w http.ResponseWriter) {

	_, ok := t[id]
	if ok {
		delete(t, id)
		fmt.Fprintf(w, "Task '%s' was deleted ", id)
	} else {
		fmt.Fprintln(w, "Task not found")
	}
}
