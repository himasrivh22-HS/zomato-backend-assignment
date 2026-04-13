package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"zomato-backend-assignment/internal/config"
	"zomato-backend-assignment/internal/handler"
)

func main() {
	// Load environment variables (optional for Docker)
	_ = godotenv.Load()

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to database")

	// 🔥 RUN MIGRATIONS
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	fmt.Println("Migrations applied")

	// Setup routes
	router := handler.SetupRoutes(db)

	// Start server
	fmt.Println("Server running on port 8080")
	fmt.Println("DB_URL:", os.Getenv("DB_URL"))

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}