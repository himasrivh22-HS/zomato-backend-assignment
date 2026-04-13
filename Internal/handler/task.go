package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"zomato-backend-assignment/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

// CREATE TASK (POST /projects/{id}/tasks)
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	// Get project ID from URL
	projectID := chi.URLParam(r, "id")

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.ProjectID = projectID // 

	query := `INSERT INTO tasks 
	(id, title, description, status, priority, project_id, assignee_id) 
	VALUES ($1,$2,$3,$4,$5,$6,$7)`

	var assignee interface{}

	if task.AssigneeID == "" {
		assignee = nil
	} else {
		assignee = task.AssigneeID
	}

	_, err = h.DB.Exec(query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ProjectID,
		assignee,
	)

	if err != nil {
		fmt.Println("TASK DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// GET TASKS (GET /projects/{id}/tasks)
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTasks API HIT")
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
	respondError(w, http.StatusBadRequest, "project id is required")
	return
}

	// Filters
	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

	//  Base query
	query := "SELECT id, title, description, status, priority, project_id, assignee_id FROM tasks WHERE project_id = $1"

	args := []interface{}{projectID}
	argIndex := 2

	if status != "" {
		query += " AND status = $" + fmt.Sprint(argIndex)
		args = append(args, status)
		argIndex++
	}

	if assignee != "" {
		query += " AND assignee_id = $" + fmt.Sprint(argIndex)
		args = append(args, assignee)
		argIndex++
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		log.Println("DB ERROR:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task
		var assigneeID sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.ProjectID,
			&assigneeID,
		)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "error reading tasks")
			return
		}

		if assigneeID.Valid {
			task.AssigneeID = assigneeID.String
		}

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}

// UPDATE TASK (PATCH /tasks/{id})
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `UPDATE tasks 
	SET title=$1, description=$2, status=$3, priority=$4 
	WHERE id=$5`

	_, err = h.DB.Exec(query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		id,
	)

	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task updated successfully"))
}

// DELETE TASK (DELETE /tasks/{id})
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM tasks WHERE id=$1"

	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task deleted successfully"))
}