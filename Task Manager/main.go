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

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
