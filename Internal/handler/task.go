package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zomato-backend-assignment/internal/middleware"
	"zomato-backend-assignment/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

// ================= CREATE TASK =================
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task model.Task

	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "project id is required")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var exists bool
	err := h.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM projects WHERE id=$1 AND owner_id=$2)",
		projectID,
		userID,
	).Scan(&exists)

	if err != nil || !exists {
		writeError(w, http.StatusForbidden, "invalid project access")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if task.Title == "" {
		writeValidationError(w, map[string]string{"title": "is required"})
		return
	}
	if task.Status == "" {
		writeValidationError(w, map[string]string{"status": "is required"})
		return
	}

	task.ID = uuid.New().String()
	task.ProjectID = projectID
	task.CreatorID = userID

	var assignee interface{}
	if task.AssigneeID == "" {
		assignee = nil
	} else {
		assignee = task.AssigneeID
	}

	var dueDate interface{}
	if task.DueDate == "" {
		dueDate = nil
	} else {
		dueDate = task.DueDate
	}

	query := `INSERT INTO tasks 
	(id, title, description, status, priority, project_id, assignee_id, creator_id, due_date, created_at, updated_at) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW())`

	_, err = h.DB.Exec(query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ProjectID,
		assignee,
		task.CreatorID,
		dueDate,
	)

	if err != nil {
		fmt.Println("TASK DB ERROR:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// ================= GET TASKS =================
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID := chi.URLParam(r, "id")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")

	// ✅ FIX: secure base query
	query := `
	SELECT id, title, description, status, priority, project_id, assignee_id 
	FROM tasks 
	WHERE project_id=$1 AND project_id IN (
		SELECT id FROM projects WHERE owner_id=$2
	)`
	args := []interface{}{projectID, userID}

	argIndex := 3

	if status != "" {
		query += " AND status=$" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if assignee != "" {
		query += " AND assignee_id=$" + strconv.Itoa(argIndex)
		args = append(args, assignee)
		argIndex++
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		var assigneeID, priority, description sql.NullString

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&description,
			&t.Status,
			&priority,
			&t.ProjectID,
			&assigneeID,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error reading tasks")
			return
		}

		if description.Valid {
			t.Description = description.String
		}
		if assigneeID.Valid {
			t.AssigneeID = assigneeID.String
		}
		if priority.Valid {
			t.Priority = priority.String
		}

		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}
// ================= UPDATE TASK =================
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	//  validation
	if task.Title == "" {
		writeValidationError(w, map[string]string{"title": "is required"})
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var projectOwnerID, assigneeID sql.NullString

	err = h.DB.QueryRow(`
		SELECT p.owner_id, t.assignee_id
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = $1
	`, id).Scan(&projectOwnerID, &assigneeID)

	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if userID != projectOwnerID.String && (!assigneeID.Valid || userID != assigneeID.String) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	query := `UPDATE tasks 
	SET title=$1, description=$2, status=$3, priority=$4, updated_at=NOW() 
	WHERE id=$5`

	_, err = h.DB.Exec(query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		id,
	)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "task updated successfully",
	})
}

// ================= DELETE TASK =================
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var projectOwnerID, creatorID sql.NullString

	err := h.DB.QueryRow(`
		SELECT p.owner_id, t.creator_id
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE t.id = $1
	`, id).Scan(&projectOwnerID, &creatorID)

	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if userID != projectOwnerID.String && (!creatorID.Valid || userID != creatorID.String) {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	_, err = h.DB.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}