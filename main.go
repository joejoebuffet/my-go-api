package main

import (
	"database/sql"
	"fmt"
	"log"
	"os" // Added to read environment variables

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Read environment variables passed from Podman/Jenkins
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")

	// 2. Set fallback defaults if they happen to be empty
	if dbHost == "" {
		dbHost = "10.36.168.15" // Your Windows PC IP
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	// 3. Dynamically build the connection string including host and port
	connStr := fmt.Sprintf(
		"host=%s port=%s user=hello_api_user password=ApiSecurePass123! dbname=mdbase sslmode=disable",
		dbHost,
		dbPort,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
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
