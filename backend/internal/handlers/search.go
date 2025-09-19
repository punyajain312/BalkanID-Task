package handlers

import (
    "encoding/json"
    "net/http"

    "balkanid-capstone/internal/services"
)

type SearchHandler struct {
    Service *services.SearchService
}

func NewSearchHandler(service *services.SearchService) *SearchHandler {
    return &SearchHandler{Service: service}
}

func (h *SearchHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok || userID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    q := r.URL.Query().Get("q")
    mime := r.URL.Query().Get("mime")
    minSize := r.URL.Query().Get("min_size")
    maxSize := r.URL.Query().Get("max_size")
    from := r.URL.Query().Get("from")
    to := r.URL.Query().Get("to")

    files, err := h.Service.SearchFiles(userID, q, mime, minSize, maxSize, from, to)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}