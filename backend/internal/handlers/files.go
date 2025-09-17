package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type FileHandler struct {
	DB *sql.DB
}

// List all files for a user (dummy user for now)
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	// TEST USER
	// userId := "11111111-1111-1111-1111-111111111111"

	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(`
        SELECT f.id, f.filename, f.mime_type, f.size, f.created_at, b.hash, b.ref_count
        FROM files f
        JOIN file_blobs b ON f.file_hash = b.hash
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
    `, userId)
	if err != nil {
		log.Println("ListFiles query error:", err)
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

// Delete a file (only decrements ref_count, deletes blob if unused)
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "missing file id", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "failed to start tx", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Find the file and its blob
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var hash string
	err = tx.QueryRow(`SELECT file_hash FROM files WHERE id=$1 userId=$2`, fileID).Scan(&hash)
	if err == sql.ErrNoRows {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete file record
	_, err = tx.Exec(`DELETE FROM files WHERE id=$1`, fileID)
	if err != nil {
		http.Error(w, "delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decrement ref_count
	_, err = tx.Exec(`UPDATE file_blobs SET ref_count = ref_count - 1 WHERE hash=$1`, hash)
	if err != nil {
		http.Error(w, "ref_count error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if blob is now unused
	var refCount int
	err = tx.QueryRow(`SELECT ref_count FROM file_blobs WHERE hash=$1`, hash).Scan(&refCount)
	if err != nil {
		http.Error(w, "ref_count check error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if refCount <= 0 {
		_, err = tx.Exec(`DELETE FROM file_blobs WHERE hash=$1`, hash)
		if err != nil {
			http.Error(w, "blob delete error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "tx commit error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
