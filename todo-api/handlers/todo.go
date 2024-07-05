package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo-api/db"
	"todo-api/models"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = gocql.UUID(uuid.New())
	todo.Created = time.Now()
	todo.Updated = time.Now()

	if err := db.Session.Query(`INSERT INTO todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")
	pageSizeStr := r.URL.Query().Get("page_size")
	pageToken := r.URL.Query().Get("page_token")
	prevPageToken := r.URL.Query().Get("prev_page_token")

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	pageSize := 10
	if pageSizeStr != "" {
		pSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
		pageSize = pSize
	}

	query := "SELECT id, user_id, title, description, status, created, updated FROM todos WHERE user_id = ?"
	params := []interface{}{userID}

	if status != "" {
		query += " AND status = ?"
		params = append(params, status)
	}

	var pagingQuery string
	var pagingToken gocql.UUID
	if pageToken != "" {
		if err := pagingToken.UnmarshalText([]byte(pageToken)); err != nil {
			http.Error(w, "Invalid page token", http.StatusBadRequest)
			return
		}
		pagingQuery = " AND token(id) > token(?)"
		params = append(params, pagingToken)
	} else if prevPageToken != "" {
		if err := pagingToken.UnmarshalText([]byte(prevPageToken)); err != nil {
			http.Error(w, "Invalid prev page token", http.StatusBadRequest)
			return
		}
		pagingQuery = " AND token(id) < token(?)"
		params = append(params, pagingToken)
	}

	query += pagingQuery + " LIMIT ?ALLOW FILTERING"
	params = append(params, pageSize)

	// Execute query
	iter := db.Session.Query(query, params...).Iter()
	defer iter.Close()

	var todos []models.Todo
	var todo models.Todo

	// Iterate through query results
	for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
		todos = append(todos, todo)
	}

	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var nextPageToken, prevPageTokenStr string
	fmt.Println(todos)
	if len(todos) > 0 {
		nextPageToken = todos[len(todos)-1].ID.String()
		if len(todos) == pageSize {
			prevPageTokenStr = todos[0].ID.String()
		}
	}
	response := map[string]interface{}{
		"todos":             todos,
		"next_page_token":   nextPageToken,
		"prev_page_token":   prevPageTokenStr,
		"current_page_size": len(todos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var todo models.Todo
	if err := db.Session.Query("SELECT id, user_id, title, description, status, created, updated FROM todos WHERE id = ?", id).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated); err != nil {
		if err == gocql.ErrNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var todoUpdate models.Todo
	json.NewDecoder(r.Body).Decode(&todoUpdate)

	var existingTodo models.Todo
	if err := db.Session.Query("SELECT id, user_id, title, description, status, created, updated FROM todos WHERE id = ?", id).Scan(
		&existingTodo.ID, &existingTodo.UserID, &existingTodo.Title, &existingTodo.Description, &existingTodo.Status, &existingTodo.Created, &existingTodo.Updated); err != nil {
		if err == gocql.ErrNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todoUpdate.Title != "" {
		existingTodo.Title = todoUpdate.Title
	}
	if todoUpdate.Description != "" {
		existingTodo.Description = todoUpdate.Description
	}
	if todoUpdate.Status != "" {
		existingTodo.Status = todoUpdate.Status
	}
	existingTodo.Updated = time.Now()

	if err := db.Session.Query(`UPDATE todos SET title = ?, description = ?, status = ?, updated = ? WHERE id = ?`,
		existingTodo.Title, existingTodo.Description, existingTodo.Status, existingTodo.Updated, id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existingTodo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if err := db.Session.Query("DELETE FROM todos WHERE id = ?", id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
