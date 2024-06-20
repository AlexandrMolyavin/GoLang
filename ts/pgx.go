package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"net/url"
)

type Pgx struct {
	db *sql.DB
}

func NewPgx() *Pgx {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   "localhost:5432",
		User:   url.UserPassword("postgres", "030200zxgta5"),
		Path:   "Tasks",
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	return &Pgx{db: db}
}

func (p *Pgx) post(w http.ResponseWriter, r *http.Request) {
	var taskDecode Task
	err := json.NewDecoder(r.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
	taskDecode.Status = "Proccessing"
	sqlReq := fmt.Sprintf("insert into task values('%s','%s')", taskDecode.Id, taskDecode.Status)
	_, err2 := p.db.ExecContext(context.Background(), sqlReq)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Task was added '%s'", taskDecode.Id)
}

func (p *Pgx) getById(id string, w http.ResponseWriter) {
	sqlReq := fmt.Sprintf("select * from task where id = '%s'", id)
	rows := p.db.QueryRow(sqlReq)

	var id_ string
	var status string

	rows.Scan(&id_, &status)
	if id_ == id {
		fmt.Fprintf(w, "Task ID: %s\nStatus: %s", id, status)
		return
	}

	fmt.Fprintln(w, "Task not found")
}

func (p *Pgx) getAll(w http.ResponseWriter) {

	rows, _ := p.db.QueryContext(context.Background(), "select * from Task")
	var tm = make(taskMap)
	var id string
	var status string
	for rows.Next() {
		rows.Scan(&id, &status)
		tm[id] = status
	}
	fmt.Fprintf(w, "%d tasks added \n", len(tm))
	if len(tm) != 0 {
		json.NewEncoder(w).Encode(tm)
	}
}

/*Изменение статуса задачи по ID*/
func (p *Pgx) changeStatus(w http.ResponseWriter, r *http.Request) {
	var taskDecode Task
	err := json.NewDecoder(r.Body).Decode(&taskDecode)
	if err != nil {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}

	sqlReq := fmt.Sprintf("update task set status = '%s' where id = '%s'", taskDecode.Status, taskDecode.Id)
	_, err2 := p.db.ExecContext(context.Background(), sqlReq)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, "Task not found")
	} else {
		fmt.Fprintf(w, "Status was updated:\nTask ID: %s\nStatus: %s", taskDecode.Id, taskDecode.Status)
	}

}

/*Удаление задачи по ID*/
func (p *Pgx) deleteById(id string, w http.ResponseWriter) {
	sqlReq := fmt.Sprintf("delete from task where id = '%s'", id)
	rows := p.db.QueryRow(sqlReq)
	if rows.Err() == nil {
		fmt.Fprintf(w, "Task '%s' was deleted ", id)
	} else {
		fmt.Fprintln(w, "Task not found")
	}
}
