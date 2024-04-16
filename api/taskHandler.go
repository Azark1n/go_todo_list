package api

import (
	"encoding/json"
	"fmt"
	"go_todo_list/errors"
	"go_todo_list/model"
	"go_todo_list/service"
	"net/http"
	"strconv"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" {
		sendErrorResponse(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	if repeatStr == "" {
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		sendErrorResponse(w, "Invalid format for 'now'", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		sendErrorResponse(w, "Invalid format for 'date'", http.StatusBadRequest)
		return
	}

	nextDate, err := service.NextDate(now, date, repeatStr)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `%s`, nextDate)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postTask(w, r)
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPut:
		putTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	id, err := service.CreateTask(task)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := parseID(w, idStr)
	if err != nil {
		return
	}

	task, err := service.GetTask(id)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, task, http.StatusOK)
}

func putTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err := service.UpdateTask(task)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, map[string]interface{}{}, http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := parseID(w, idStr)
	if err != nil {
		return
	}

	err = service.DeleteTask(id)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, map[string]interface{}{}, http.StatusOK)
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := parseID(w, idStr)
	if err != nil {
		return
	}

	err = service.MarkDoneTask(id)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, map[string]interface{}{}, http.StatusOK)
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")
	tasks, err := service.FindTasks(search)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	sendJSON(w, map[string][]model.Task{"tasks": tasks}, http.StatusOK)
}

func handleErrorResponse(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *errors.ValidationError:
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
	case *errors.NotFoundError:
		sendErrorResponse(w, err.Error(), http.StatusNotFound)
	default:
		sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
	}
}

func parseID(w http.ResponseWriter, idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendErrorResponse(w, "id is not correct", http.StatusBadRequest)
		return 0, err
	}
	return id, nil
}

func sendErrorResponse(w http.ResponseWriter, err string, code int) {
	sendJSON(w, map[string]string{"errors": err}, code)
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
