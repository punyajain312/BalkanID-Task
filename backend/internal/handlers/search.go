package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"
)

type SearchHandler struct {
    DB *sql.DB
}

func (h *SearchHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
    userId, ok := r.Context().Value("user_id").(string)
    if !ok || userId == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // Query params
    q := r.URL.Query().Get("q")
    mime := r.URL.Query().Get("mime")
    minSize := r.URL.Query().Get("min_size")
    maxSize := r.URL.Query().Get("max_size")
    from := r.URL.Query().Get("from")
    to := r.URL.Query().Get("to")

    // Base query
    query := `
        SELECT f.id, f.filename, f.mime_type, f.size, f.created_at, b.hash, b.ref_count
        FROM files f
        JOIN file_blobs b ON f.blob_id = b.id
        WHERE f.user_id = $1
    `
    args := []interface{}{userId}
    argIndex := 2

    // Dynamic filters
    if q != "" {
        query += " AND f.filename ILIKE $" + strconv.Itoa(argIndex)
        args = append(args, "%"+q+"%")
        argIndex++
    }
    if mime != "" {
        query += " AND f.mime_type = $" + strconv.Itoa(argIndex)
        args = append(args, mime)
        argIndex++
    }
    if minSize != "" {
        if val, err := strconv.ParseInt(minSize, 10, 64); err == nil {
            query += " AND f.size >= $" + strconv.Itoa(argIndex)
            args = append(args, val)
            argIndex++
        }
    }
    if maxSize != "" {
        if val, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
            query += " AND f.size <= $" + strconv.Itoa(argIndex)
            args = append(args, val)
            argIndex++
        }
    }
    if from != "" {
        if t, err := time.Parse("2006-01-02", from); err == nil {
            query += " AND f.created_at >= $" + strconv.Itoa(argIndex)
            args = append(args, t)
            argIndex++
        }
    }
    if to != "" {
        if t, err := time.Parse("2006-01-02", to); err == nil {
            query += " AND f.created_at <= $" + strconv.Itoa(argIndex)
            args = append(args, t)
            argIndex++
        }
    }

    query += " ORDER BY f.created_at DESC"

    rows, err := h.DB.Query(query, args...)
    if err != nil {
        log.Println("Search query error:", err)
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type File struct {
        ID        string `json:"id"`
        Filename  string `json:"filename"`
        MimeType  string `json:"mime_type"`
        Size      int64  `json:"size"`
        CreatedAt string `json:"created_at"`
        Hash      string `json:"hash"`
        RefCount  int    `json:"ref_count"`
    }

    files := []File{}
    for rows.Next() {
        var f File
        if err := rows.Scan(&f.ID, &f.Filename, &f.MimeType, &f.Size, &f.CreatedAt, &f.Hash, &f.RefCount); err != nil {
            http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        files = append(files, f)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}