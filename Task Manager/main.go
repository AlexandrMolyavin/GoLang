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

var taskServ []Task

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
			tasksRep.getAll(w)
		} else {
			tasksRep.getById(r.URL.Path[1:], w)
		}
	case "POST":
		tasksRep.post(w, r)
	case "DELETE":
		tasksRep.deleteById(r.URL.Path[1:], w)
	case "PUT":
		tasksRep.changeStatus(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}
func getTask(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%d tasks added \n", len(taskServ))
	if len(taskServ) != 0 {
		json.NewEncoder(w).Encode(taskServ)
	}

}
func postTask(w http.ResponseWriter, r *http.Request) {
	var taskDecode Task
	err := json.NewDecoder(r.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
	taskDecode.Status = "Proccessing"
	taskServ = append(taskServ, taskDecode)
	fmt.Fprintf(w, "Task was added '%s'", taskDecode.Id)
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	if len(taskServ) != 0 {
		taskServ = append(taskServ[:0], taskServ[1:]...)
		fmt.Fprintln(w, "First task was deleted")
	} else {
		fmt.Fprintln(w, "All tasks were deleted ")
	}
}

@func main() {
	//tasks = json_reader()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
