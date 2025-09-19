package main

import (
	"fmt"
	"log"
	"net/http"

	"balkanid-capstone/internal/db"
	"balkanid-capstone/internal/handlers"
	"balkanid-capstone/internal/middleware"

	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to Postgres
	dsn := "postgres://postgres:1234@localhost:5432/balkanid?sslmode=disable"
	database, err := db.Connect(dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.Close()

	// Handlers
	uploadHandler := &handlers.UploadHandler{DB: database}
	fileHandler := &handlers.FileHandler{DB: database}
	searchHandler := &handlers.SearchHandler{DB: database}
	statsHandler := &handlers.StatsHandler{DB: database}
	authHandler := &handlers.AuthHandler{DB: database}
	adminHandler := &handlers.AdminHandler{DB: database}

	// Router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)

	// Protected routes
	mux.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(fileHandler.ListFiles)))
	mux.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(uploadHandler.UploadFile)))
	mux.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(fileHandler.DeleteFile)))
	mux.Handle("/search", middleware.AuthMiddleware(http.HandlerFunc(searchHandler.SearchFiles)))
	mux.Handle("/stats", middleware.AuthMiddleware(http.HandlerFunc(statsHandler.GetStats)))

	// Admin routes
	mux.Handle("/admin/users", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListUsers))))
	mux.Handle("/admin/files", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListAllFiles))))
	mux.Handle("/admin/stats", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.SystemStats))))

	// Add CORS so frontend can talk to backend
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}