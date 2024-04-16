package main

import (
	"go_todo_list/api"
	"go_todo_list/config"
	"go_todo_list/data"
	"log"
	"net/http"
)

func main() {
	config.AcceptEnvironments()

	data.OpenDbOrCreate()
	defer data.Db.Close()

	http.Handle("/", http.FileServer(http.Dir(config.WebDir)))
	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.Auth(api.TaskHandler))
	http.HandleFunc("/api/task/done", api.Auth(api.TaskDoneHandler))
	http.HandleFunc("/api/tasks", api.Auth(api.TasksHandler))
	http.HandleFunc("/api/signin", api.SignInHandler)

	log.Printf("Starting server on :%s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
