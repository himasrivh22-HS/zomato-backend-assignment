package main

import (
	"fmt"
	"net/http"
	"zomato-backend-assignment/internal/handler"
	"zomato-backend-assignment/internal/config"
	"zomato-backend-assignment/internal/middleware"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	authHandler := &handler.AuthHandler{DB: db}
	projectHandler := &handler.ProjectHandler{DB: db}
	taskHandler := &handler.TaskHandler{DB: db} 

	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)

	http.HandleFunc("/projects", projectHandler.CreateProject)

	http.HandleFunc("/tasks", middleware.AuthMiddleware(taskHandler.CreateTask))
    http.HandleFunc("/tasks/list", middleware.AuthMiddleware(taskHandler.GetTasks))
    http.HandleFunc("/tasks/update", middleware.AuthMiddleware(taskHandler.UpdateTask))
    http.HandleFunc("/tasks/delete", middleware.AuthMiddleware(taskHandler.DeleteTask))

	fmt.Println("Connected to database ")
	fmt.Println("Server running on port 8080 ")

	router := handler.SetupRoutes(db)
    http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}