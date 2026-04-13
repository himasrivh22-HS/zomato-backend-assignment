package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"zomato-backend-assignment/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	DB *sql.DB
}

type ProjectResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	OwnerID     string       `json:"owner_id"`
	Tasks       []model.Task `json:"tasks"`
}

// ================= CREATE PROJECT =================
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	project.ID = uuid.New().String()
	project.OwnerID = userID

	query := `INSERT INTO projects (id, name, description, owner_id)
	          VALUES ($1, $2, $3, $4)`

	_, err = h.DB.Exec(query, project.ID, project.Name, project.Description, project.OwnerID)
	if err != nil {
		fmt.Println("PROJECT DB ERROR:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
    w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// ================= GET ALL PROJECTS =================
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, description, owner_id FROM projects")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}
	defer rows.Close()

	var projects []model.Project

	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error reading projects")
			return
		}
		projects = append(projects, p)
	}

	json.NewEncoder(w).Encode(projects)
}

// ================= GET PROJECT BY ID =================
func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p model.Project

	err := h.DB.QueryRow("SELECT id, name, description, owner_id FROM projects WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)

	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, title, description, status, priority, project_id, assignee_id 
		FROM tasks WHERE project_id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		var desc, priority, assignee sql.NullString

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&desc,
			&t.Status,
			&priority,
			&t.ProjectID,
			&assignee,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error reading tasks")
			return
		}

		if desc.Valid {
			t.Description = desc.String
		}
		if priority.Valid {
			t.Priority = priority.String
		}
		if assignee.Valid {
			t.AssigneeID = assignee.String
		}

		tasks = append(tasks, t)
	}

	response := ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		Tasks:       tasks,
	}

	json.NewEncoder(w).Encode(response)
}

// ================= UPDATE PROJECT =================
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p model.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var ownerID string
	err := h.DB.QueryRow("SELECT owner_id FROM projects WHERE id=$1", id).Scan(&ownerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	if ownerID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	_, err = h.DB.Exec("UPDATE projects SET name=$1, description=$2 WHERE id=$3",
		p.Name, p.Description, id)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	w.Write([]byte("Project updated"))
}

// ================= DELETE PROJECT =================
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var ownerID string
	err := h.DB.QueryRow("SELECT owner_id FROM projects WHERE id=$1", id).Scan(&ownerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	if ownerID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	_, err = h.DB.Exec("DELETE FROM projects WHERE id=$1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}