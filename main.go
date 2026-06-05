package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Added for local .env handling
	_ "github.com/lib/pq"
)

func main() {
	// 1. Try to load local .env file. If it doesn't exist (like on your Podman node),
	// it fails silently and logs an info message without crashing.
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: No .env file found. Relying on system environment variables.")
	}

	// 2. Read environment variables passed from Podman/Jenkins/Local .env
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DB_CREDS_USR")
	dbPass := os.Getenv("DB_CREDS_PSW")
	// 3. Set fallback defaults if they happen to be empty
	if dbHost == "" {
		dbHost = "10.36.168.15" // Your Windows PC IP
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	// 4. Dynamically build the connection string including host and port
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=mdbase sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPass,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to initialize postgres database: %v", err)
	}
	defer db.Close()

	// Instantiate the controller with its dependency injection
	helloCtrl := NewHelloController(db)

	router := gin.Default()

	// Direct matching mappings
	router.GET("/hello", helloCtrl.GetAccountHolder)
	router.POST("/updateStatus", helloCtrl.UpdateAccountStatus)
	router.POST("/updateStatusBulk", helloCtrl.UpdateAccountStatusBulk)

	router.Run(":8081")
}
