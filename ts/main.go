package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type manager interface {
	post(w http.ResponseWriter, r *http.Request)
	getById(id string, w http.ResponseWriter)
	getAll(w http.ResponseWriter)
	changeStatus(w http.ResponseWriter, r *http.Request)
	deleteById(id string, w http.ResponseWriter)
}

/*Получения задач из .json файла*/
func json_reader() []Task {

	data, err := ioutil.ReadFile("E:\\go_projects\\Education\\Task_Handler\\Tasks.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var task []Task

	err = json.Unmarshal(data, &task)
	if err != nil {
		fmt.Println("Error JSON unmarshalling tasks.")
		fmt.Println(err.Error())
		return nil
	}

	return task
}

func handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		if len(r.URL.Path[1:]) == 0 {
			//tasksRep.getAll(w)
			pg.getAll(w)
		} else {
			//tasksRep.getById(r.URL.Path[1:], w)
			pg.getById(r.URL.Path[1:], w)
		}
	case "POST":
		//tasksRep.post(w, r)
		pg.post(w, r)
	case "DELETE":
		//tasksRep.deleteById(r.URL.Path[1:], w)
		pg.deleteById(r.URL.Path[1:], w)
	case "PUT":
		//tasksRep.changeStatus(w, r)
		pg.changeStatus(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
