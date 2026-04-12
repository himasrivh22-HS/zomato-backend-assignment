package handler

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"net/http"

	"zomato-backend-assignment/internal/model"

	"github.com/google/uuid"
)

type ProjectHandler struct {
	DB *sql.DB
}

// CREATE PROJECT
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	project.ID = uuid.New().String()
	project.OwnerID = uuid.New().String()

	query := "INSERT INTO projects (id, name, description, owner_id) VALUES ($1,$2,$3,$4)"

	_, err = h.DB.Exec(query, project.ID, project.Name, project.Description, project.OwnerID)
if err != nil {
	fmt.Println("PROJECT DB ERROR:", err)  
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

	json.NewEncoder(w).Encode(project)
}