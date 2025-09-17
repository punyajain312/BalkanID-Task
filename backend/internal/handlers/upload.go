package handlers

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "os"
)

type UploadHandler struct {
    DB *sql.DB
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // 1. Parse multipart file
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "failed to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Read file into memory (to hash + save later)
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "failed to read file", http.StatusInternalServerError)
        return
    }

    // 2. Compute SHA-256
    hash := sha256.Sum256(fileBytes)
    hashStr := hex.EncodeToString(hash[:])

    // 3. Check if blob exists
    var exists bool
    err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM file_blobs WHERE hash=$1)", hashStr).Scan(&exists)
    if err != nil {
        http.Error(w, "db error", http.StatusInternalServerError)
        return
    }

    if exists {
        // Increment ref_count
        _, _ = h.DB.Exec("UPDATE file_blobs SET ref_count = ref_count + 1 WHERE hash=$1", hashStr)
    } else {
        // Save file to local storage
        path := "uploads/" + hashStr
        err = os.WriteFile(path, fileBytes, 0644)
        if err != nil {
            http.Error(w, "failed to save file", http.StatusInternalServerError)
            return
        }

        // Insert into file_blobs
        _, _ = h.DB.Exec(`
            INSERT INTO file_blobs (hash, storage_path, size, ref_count) 
            VALUES ($1, $2, $3, 1)`,
            hashStr, path, header.Size)
    }

	// Test User ID
    // userId := "11111111-1111-1111-1111-111111111111"

    // 4. Insert into files table
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_, _ = h.DB.Exec(`
		INSERT INTO files (user_id, file_hash, filename, mime_type, size) 
		VALUES ($1, $2, $3, $4, $5)`,
		userId, hashStr, header.Filename, header.Header.Get("Content-Type"), header.Size)

    // 5. Respond
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "hash":     hashStr,
        "filename": header.Filename,
    })
}
