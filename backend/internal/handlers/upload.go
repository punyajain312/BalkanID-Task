package handlers

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type UploadHandler struct {
    DB *sql.DB
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // 1. Parse multipart form (max 20MB)
    if err := r.ParseMultipartForm(20 << 20); err != nil {
        http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
        return
    }

    files := r.MultipartForm.File["files"]
    if len(files) == 0 {
        http.Error(w, "no files uploaded", http.StatusBadRequest)
        return
    }

    // Get user ID from JWT context
    userId, ok := r.Context().Value("user_id").(string)
    if !ok || userId == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    uploaded := []map[string]string{}

    for _, header := range files {
        file, err := header.Open()
        if err != nil {
            http.Error(w, "failed to open file", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        fileBytes, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "failed to read file", http.StatusInternalServerError)
            return
        }

        // 2. Compute SHA-256
        hash := sha256.Sum256(fileBytes)
        hashStr := hex.EncodeToString(hash[:])

        // 3. Look for existing blob by hash
        var blobID string
        var exists bool
        err = h.DB.QueryRow(`
            SELECT id, true FROM file_blobs WHERE hash=$1
        `, hashStr).Scan(&blobID, &exists)

        if err == sql.ErrNoRows {
            // Blob does not exist → insert new
            os.MkdirAll("uploads", os.ModePerm)
            path := filepath.Join("uploads", hashStr)

            if err := os.WriteFile(path, fileBytes, 0644); err != nil {
                http.Error(w, "failed to save file", http.StatusInternalServerError)
                return
            }

            err = h.DB.QueryRow(`
                INSERT INTO file_blobs (hash, storage_path, size, ref_count)
                VALUES ($1, $2, $3, 1) RETURNING id
            `, hashStr, path, header.Size).Scan(&blobID)
            if err != nil {
                http.Error(w, "db insert error: "+err.Error(), http.StatusInternalServerError)
                return
            }
        } else if err != nil {
            http.Error(w, "db query error: "+err.Error(), http.StatusInternalServerError)
            return
        } else {
            // Blob exists → increment ref_count
            _, _ = h.DB.Exec("UPDATE file_blobs SET ref_count = ref_count + 1 WHERE id=$1", blobID)
        }

        // 4. Insert into files table
        _, err = h.DB.Exec(`
            INSERT INTO files (user_id, blob_id, filename, mime_type, size)
            VALUES ($1, $2, $3, $4, $5)
        `, userId, blobID, header.Filename, header.Header.Get("Content-Type"), header.Size)
        if err != nil {
            http.Error(w, "db insert error: "+err.Error(), http.StatusInternalServerError)
            return
        }

        uploaded = append(uploaded, map[string]string{
            "hash":     hashStr,
            "filename": header.Filename,
        })
    }

    // 5. Respond
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "upload successful",
        "files":   uploaded,
    })
}