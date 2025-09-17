package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
)

type SearchHandler struct {
    DB *sql.DB
}

func (h *SearchHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {

	// TEST USER
    // userId := "11111111-1111-1111-1111-111111111111" 

	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

    // Query params
    q := r.URL.Query().Get("q")            // filename search
    mime := r.URL.Query().Get("mime")      // MIME type filter
    minSize := r.URL.Query().Get("min_size")
    maxSize := r.URL.Query().Get("max_size")
    from := r.URL.Query().Get("from")      // date range
    to := r.URL.Query().Get("to")

    // Base query
    query := `
        SELECT f.id, f.filename, f.mime_type, f.size, f.created_at, b.hash, b.ref_count
        FROM files f
        JOIN file_blobs b ON f.file_hash = b.hash
        WHERE f.user_id = $1
    `
    args := []interface{}{userId}
    argIndex := 2

    // Dynamic filters
    if q != "" {
        query += " AND f.filename ILIKE $" + string(rune(argIndex))
        args = append(args, "%"+q+"%")
        argIndex++
    }
    if mime != "" {
        query += " AND f.mime_type = $" + string(rune(argIndex))
        args = append(args, mime)
        argIndex++
    }
    if minSize != "" {
        query += " AND f.size >= $" + string(rune(argIndex))
        args = append(args, minSize)
        argIndex++
    }
    if maxSize != "" {
        query += " AND f.size <= $" + string(rune(argIndex))
        args = append(args, maxSize)
        argIndex++
    }
    if from != "" {
        query += " AND f.created_at >= $" + string(rune(argIndex))
        args = append(args, from)
        argIndex++
    }
    if to != "" {
        query += " AND f.created_at <= $" + string(rune(argIndex))
        args = append(args, to)
        argIndex++
    }

    query += " ORDER BY f.created_at DESC"

    // Run query
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

    json.NewEncoder(w).Encode(files)
}
