package handlers

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"

    "balkanid-capstone/internal/db"
    "balkanid-capstone/internal/middleware"

    _ "github.com/lib/pq"
)

var testToken string
var testAdminToken string

// Helper: create multipart file upload body
func createMultipartFile(t *testing.T, field, filename, content string) (*bytes.Buffer, string) {
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(field, filename)
    if err != nil {
        t.Fatal(err)
    }
    io.Copy(part, strings.NewReader(content))
    writer.Close()
    return body, writer.FormDataContentType()
}

// Setup test server with DB + routes
func setupTestServer(t *testing.T) (*httptest.Server, *sql.DB) {
    dsn := os.Getenv("TEST_DSN") // e.g. postgres://postgres:pass@localhost:5432/balkanid_test?sslmode=disable
	if dsn == "" {
		dsn = "postgres://postgres:yourpassword@localhost:5432/balkanid_test?sslmode=disable"
	}

    dbConn, err := db.Connect(dsn)
    if err != nil {
        t.Fatal("DB connect failed:", err)
    }

    // Handlers
    authHandler := &AuthHandler{DB: dbConn}
    uploadHandler := &UploadHandler{DB: dbConn}
    fileHandler := &FileHandler{DB: dbConn}
    searchHandler := &SearchHandler{DB: dbConn}
    statsHandler := &StatsHandler{DB: dbConn}
    adminHandler := &AdminHandler{DB: dbConn}

    mux := http.NewServeMux()

    // Public
    mux.HandleFunc("/signup", authHandler.Signup)
    mux.HandleFunc("/login", authHandler.Login)

    // Protected user routes
    mux.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(uploadHandler.UploadFile)))
    mux.Handle("/files", middleware.AuthMiddleware(http.HandlerFunc(fileHandler.ListFiles)))
    mux.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(fileHandler.DeleteFile)))
    mux.Handle("/search", middleware.AuthMiddleware(http.HandlerFunc(searchHandler.SearchFiles)))
    mux.Handle("/stats", middleware.AuthMiddleware(http.HandlerFunc(statsHandler.GetStats)))

    // Admin routes
    mux.Handle("/admin/users", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListUsers))))
    mux.Handle("/admin/files", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.ListAllFiles))))
    mux.Handle("/admin/stats", middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(adminHandler.SystemStats))))

    return httptest.NewServer(mux), dbConn
}

func TestAPIFlow(t *testing.T) {
    server, dbConn := setupTestServer(t)
    defer server.Close()
    defer dbConn.Close()

    // --- Signup user ---
    resp, err := http.Post(server.URL+"/signup", "application/json",
        strings.NewReader(`{"email":"user@test.com","password":"secret"}`))
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("signup failed: %v", resp.Status)
    }

    // --- Login user ---
    resp, err = http.Post(server.URL+"/login", "application/json",
        strings.NewReader(`{"email":"user@test.com","password":"secret"}`))
    if err != nil {
        t.Fatal(err)
    }
    var login map[string]string
    json.NewDecoder(resp.Body).Decode(&login)
    testToken = login["token"]

    // --- Upload file ---
    body, contentType := createMultipartFile(t, "file", "test.txt", "hello world")
    req, _ := http.NewRequest("POST", server.URL+"/upload", body)
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Authorization", "Bearer "+testToken)
    resp, err = http.DefaultClient.Do(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("upload failed: %v", resp.Status)
    }

    // --- List files ---
    req, _ = http.NewRequest("GET", server.URL+"/files", nil)
    req.Header.Set("Authorization", "Bearer "+testToken)
    resp, err = http.DefaultClient.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        t.Fatalf("list files failed: %v", resp.Status)
    }

    // --- Stats ---
    req, _ = http.NewRequest("GET", server.URL+"/stats", nil)
    req.Header.Set("Authorization", "Bearer "+testToken)
    resp, err = http.DefaultClient.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        t.Fatalf("stats failed: %v", resp.Status)
    }

    // --- Signup admin ---
    _, err = http.Post(server.URL+"/signup", "application/json",
        strings.NewReader(`{"email":"admin@test.com","password":"adminpass"}`))
    if err != nil {
        t.Fatal(err)
    }

    // Promote to admin
    _, err = dbConn.Exec(`UPDATE users SET role='admin' WHERE email='admin@test.com'`)
    if err != nil {
        t.Fatal("failed to promote admin:", err)
    }

    // --- Admin login ---
    resp, err = http.Post(server.URL+"/login", "application/json",
        strings.NewReader(`{"email":"admin@test.com","password":"adminpass"}`))
    if err != nil {
        t.Fatal(err)
    }
    var adminLogin map[string]string
    json.NewDecoder(resp.Body).Decode(&adminLogin)
    testAdminToken = adminLogin["token"]

    // --- Admin stats ---
    req, _ = http.NewRequest("GET", server.URL+"/admin/stats", nil)
    req.Header.Set("Authorization", "Bearer "+testAdminToken)
    resp, err = http.DefaultClient.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        t.Fatalf("admin stats failed: %v", resp.Status)
    }
}
