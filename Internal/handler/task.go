package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"zomato-backend-assignment/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

// ================= CREATE TASK =================
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "project id is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.ProjectID = projectID

	var assignee interface{}
	if task.AssigneeID == "" {
		assignee = nil
	} else {
		assignee = task.AssigneeID
	}

	query := `INSERT INTO tasks 
	(id, title, description, status, priority, project_id, assignee_id) 
	VALUES ($1,$2,$3,$4,$5,$6,$7)`

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// ================= GET TASKS =================
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("GetTasks API HIT")

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "project id is required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

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
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
			http.Error(w, "error reading tasks", http.StatusInternalServerError)
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

// ================= UPDATE TASK =================
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//  Get user safely
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//  Check ownership
	var projectOwnerID, assigneeID sql.NullString

	err = h.DB.QueryRow(`
		SELECT p.owner_id, t.assignee_id
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = $1
	`, id).Scan(&projectOwnerID, &assigneeID)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if userID != projectOwnerID.String && (!assigneeID.Valid || userID != assigneeID.String) {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

// ================= DELETE TASK =================
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	//  Get user safely
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//  Check ownership
	var projectOwnerID, assigneeID sql.NullString

	err := h.DB.QueryRow(`
		SELECT p.owner_id, t.assignee_id
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = $1
	`, id).Scan(&projectOwnerID, &assigneeID)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if userID != projectOwnerID.String && (!assigneeID.Valid || userID != assigneeID.String) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	query := "DELETE FROM tasks WHERE id=$1"
	_, err = h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task deleted successfully"))
}