package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

type AdminHandler struct {
    DB *sql.DB
}

// List all users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := h.DB.Query(`SELECT id, email, role, created_at FROM users`)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type User struct {
        ID        string `json:"id"`
        Email     string `json:"email"`
        Role      string `json:"role"`
        CreatedAt string `json:"created_at"`
    }
    users := []User{}
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt); err != nil {
            http.Error(w, "scan error", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }
    json.NewEncoder(w).Encode(users)
}

// List all files
func (h *AdminHandler) ListAllFiles(w http.ResponseWriter, r *http.Request) {
    rows, err := h.DB.Query(`
        SELECT f.id, f.filename, f.size, f.created_at, u.email
        FROM files f
        JOIN users u ON f.user_id = u.id
        ORDER BY f.created_at DESC
    `)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type File struct {
        ID        string `json:"id"`
        Filename  string `json:"filename"`
        Size      int64  `json:"size"`
        CreatedAt string `json:"created_at"`
        Owner     string `json:"owner_email"`
    }
    files := []File{}
    for rows.Next() {
        var f File
        if err := rows.Scan(&f.ID, &f.Filename, &f.Size, &f.CreatedAt, &f.Owner); err != nil {
            http.Error(w, "scan error", http.StatusInternalServerError)
            return
        }
        files = append(files, f)
    }
    json.NewEncoder(w).Encode(files)
}

// System-wide storage stats
func (h *AdminHandler) SystemStats(w http.ResponseWriter, r *http.Request) {
    var totalUsers, totalFiles int
    var totalSize, uniqueSize int64

    h.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
    h.DB.QueryRow(`SELECT COUNT(*) FROM files`).Scan(&totalFiles)
    h.DB.QueryRow(`SELECT COALESCE(SUM(size),0) FROM files`).Scan(&totalSize)
    h.DB.QueryRow(`SELECT COALESCE(SUM(size),0) FROM file_blobs`).Scan(&uniqueSize)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "total_users": totalUsers,
        "total_files": totalFiles,
        "total_size": totalSize,
        "unique_size": uniqueSize,
        "savings": totalSize - uniqueSize,
    })
}
