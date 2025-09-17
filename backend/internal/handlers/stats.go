package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

type StatsHandler struct {
    DB *sql.DB
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// TEST USER
    // userId := "11111111-1111-1111-1111-111111111111" 

	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

    var totalFiles int
    var totalSize int64

    // Count total files and size for this user
    err := h.DB.QueryRow(`
        SELECT COUNT(*), COALESCE(SUM(size), 0)
        FROM files
        WHERE user_id = $1
    `, userId).Scan(&totalFiles, &totalSize)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Unique size = sum of blobs this user references
    var uniqueSize int64
    err = h.DB.QueryRow(`
        SELECT COALESCE(SUM(DISTINCT b.size), 0)
        FROM files f
        JOIN file_blobs b ON f.file_hash = b.hash
        WHERE f.user_id = $1
    `, userId).Scan(&uniqueSize)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    savings := totalSize - uniqueSize
    savingsPercent := 0.0
    if totalSize > 0 {
        savingsPercent = (float64(savings) / float64(totalSize)) * 100
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "total_files":     totalFiles,
        "total_size":      totalSize,
        "unique_size":     uniqueSize,
        "savings":         savings,
        "savings_percent": savingsPercent,
    })
}
