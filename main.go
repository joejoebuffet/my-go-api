package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Replace with your target database driver (e.g., github.com/godror/godror for Oracle)
)

func main() {
	// Configure connection string (Normally dynamically injected via OS Environment variables)
	connStr := "user=hello_api_user password=ApiSecurePass123! dbname=mdbase sslmode=disable"

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
