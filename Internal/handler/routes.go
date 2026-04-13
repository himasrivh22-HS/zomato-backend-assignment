package handler

import (
	"database/sql"
	"net/http"
	"zomato-backend-assignment/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	taskHandler := &TaskHandler{DB: db}
	projectHandler := &ProjectHandler{DB: db}
	authHandler := &AuthHandler{DB: db}

	// -------- AUTH ROUTES (NO TOKEN REQUIRED) --------
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// -------- PROTECTED ROUTES --------
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// PROJECT ROUTES
		r.Get("/projects", projectHandler.GetProjects)
		r.Post("/projects", projectHandler.CreateProject)
		r.Get("/projects/{id}", projectHandler.GetProjectByID)
		r.Patch("/projects/{id}", projectHandler.UpdateProject)
		r.Delete("/projects/{id}", projectHandler.DeleteProject)

		// TASK ROUTES
		r.Get("/projects/{id}/tasks", taskHandler.GetTasks)
		r.Post("/projects/{id}/tasks", taskHandler.CreateTask)
		r.Patch("/tasks/{id}", taskHandler.UpdateTask)
		r.Delete("/tasks/{id}", taskHandler.DeleteTask)
	})

	return r
}