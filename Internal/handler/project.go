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

// ================= CREATE PROJECT =================
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//  Get logged-in user safely
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	project.ID = uuid.New().String()
	project.OwnerID = userID 

	query := `INSERT INTO projects (id, name, description, owner_id)
	          VALUES ($1, $2, $3, $4)`

	_, err = h.DB.Exec(query,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
	)

	if err != nil {
		fmt.Println("PROJECT DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// ================= GET ALL PROJECTS =================
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, description, owner_id FROM projects")
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []model.Project

	for rows.Next() {
		var p model.Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
		if err != nil {
			http.Error(w, "Error reading projects", http.StatusInternalServerError)
			return
		}
		projects = append(projects, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// ================= GET PROJECT BY ID =================
func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p model.Project

	query := "SELECT id, name, description, owner_id FROM projects WHERE id=$1"

	err := h.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.OwnerID,
	)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// ================= UPDATE PROJECT =================
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p model.Project
	err := json.NewDecoder(r.Body).Decode(&p)
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
	var ownerID string
	err = h.DB.QueryRow("SELECT owner_id FROM projects WHERE id=$1", id).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	//  Update
	query := `UPDATE projects SET name=$1, description=$2 WHERE id=$3`

	_, err = h.DB.Exec(query, p.Name, p.Description, id)
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Project updated"))
}

// ================= DELETE PROJECT =================
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	//  Get user safely
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//  Check ownership
	var ownerID string
	err := h.DB.QueryRow("SELECT owner_id FROM projects WHERE id=$1", id).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	//  Delete
	query := "DELETE FROM projects WHERE id=$1"
	_, err = h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Project deleted"))
}