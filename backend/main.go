package main

import (
	"fmt"
	"log"
	"net/http"

	"balkanid-capstone/internal/db"
	"balkanid-capstone/internal/handlers"
	"balkanid-capstone/internal/middleware"

	_ "github.com/lib/pq" // Postgres driver
)

func main() {
	// Connect to Postgres
	dsn := "postgres://postgres:1234@localhost:5432/balkanid?sslmode=disable"
	database, err := db.Connect(dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()

	// Upload handler with DB injected
	uploadHandler := &handlers.UploadHandler{DB: database}

	// File handler
	fileHandler := &handlers.FileHandler{DB: database}

	// Search handler
	searchHandler := &handlers.SearchHandler{DB: database}

	// Stats handler
	statsHandler := &handlers.StatsHandler{DB: database}

	// Auth handler
	authHandler := &handlers.AuthHandler{DB: database}

	// Admin handler
	adminHandler := &handlers.AdminHandler{DB: database}

	// Public
    http.HandleFunc("/signup", authHandler.Signup)
    http.HandleFunc("/login", authHandler.Login)

	// Protected Routes
	http.HandleFunc("/upload", uploadHandler.UploadFile) // POST
	http.HandleFunc("/files", fileHandler.ListFiles)     // GET
	http.HandleFunc("/delete", fileHandler.DeleteFile)	// DELETE
	http.HandleFunc("/search", searchHandler.SearchFiles) // SEARCH
	http.HandleFunc("/stats", statsHandler.GetStats)	// STATS

	// Admin-Only Routes
	http.Handle("/admin/users", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListUsers))))
    http.Handle("/admin/files", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListAllFiles))))
    http.Handle("/admin/stats", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.SystemStats))))

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

