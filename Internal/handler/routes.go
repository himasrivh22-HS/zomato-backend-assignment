package handler

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// Create handler instances
	taskHandler := &TaskHandler{DB: db}
	projectHandler := &ProjectHandler{DB: db}
	authHandler := &AuthHandler{DB: db} // 👈 ADD THIS

	// -------- AUTH ROUTES --------
	r.Post("/auth/register", authHandler.Register) // 👈 FIX
	r.Post("/auth/login", authHandler.Login)       // 👈 FIX

	// -------- PROJECT ROUTES --------
	r.Get("/projects", projectHandler.GetProjects)
	r.Post("/projects", projectHandler.CreateProject)
	r.Get("/projects/{id}", projectHandler.GetProjectByID)
	r.Patch("/projects/{id}", projectHandler.UpdateProject)
	r.Delete("/projects/{id}", projectHandler.DeleteProject)

	// -------- TASK ROUTES --------
	r.Get("/projects/{id}/tasks", taskHandler.GetTasks)
	r.Post("/projects/{id}/tasks", taskHandler.CreateTask)
	r.Patch("/tasks/{id}", taskHandler.UpdateTask)
	r.Delete("/tasks/{id}", taskHandler.DeleteTask)

	return r
}