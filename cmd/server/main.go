package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"zomato-backend-assignment/internal/config"
	"zomato-backend-assignment/internal/handler"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found (continuing...)")
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to database")
	fmt.Println("Server running on port 8080")
	fmt.Println("DB_URL:", os.Getenv("DB_URL"))

	// Setup routes using chi router
	router := handler.SetupRoutes(db)

	// Start server
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}