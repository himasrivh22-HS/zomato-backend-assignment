package main

import (
	"fmt"
	"net/http"

	"zomato-backend-assignment/internal/config"
	"zomato-backend-assignment/internal/handler"
)

func main() {
	db, err := config.ConnectDB()
if err != nil {
	panic(err)
}
defer db.Close()

authHandler := &handler.AuthHandler{DB: db}

http.HandleFunc("/auth/register", authHandler.Register)
http.HandleFunc("/auth/login", authHandler.Login)

fmt.Println("Server starting on port 8080...")
err = http.ListenAndServe(":8080", nil)
if err != nil {
	panic(err)
}
}

