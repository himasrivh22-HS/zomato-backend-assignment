package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"zomato-backend-assignment/internal/model"

	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

// CREATE TASK
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()

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

// GET TASKS
// GET TASKS
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, description, status, priority, project_id, assignee_id FROM tasks")
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task
		var assignee sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.ProjectID,
			&assignee,
		)
		if err != nil {
			http.Error(w, "Error reading tasks", http.StatusInternalServerError)
			return
		}

		if assignee.Valid {
			task.AssigneeID = assignee.String
		} else {
			task.AssigneeID = ""
		}

		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

// UPDATE TASK
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4 WHERE id=$5`

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

	w.Write([]byte("Task updated successfully ✅"))
}

// DELETE TASK
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	query := "DELETE FROM tasks WHERE id=$1"

	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task deleted successfully "))
}