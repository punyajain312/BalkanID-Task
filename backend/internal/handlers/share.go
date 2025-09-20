package handlers

import (
    "encoding/json"
    "net/http"

    "balkanid-capstone/internal/models"
    "balkanid-capstone/internal/services"
)

type ShareHandler struct {
    Service *services.ShareService
}

func NewShareHandler(service *services.ShareService) *ShareHandler {
    return &ShareHandler{Service: service}
}

// POST
func (h *ShareHandler) CreateShare(w http.ResponseWriter, r *http.Request) {
    ownerID, ok := r.Context().Value("user_id").(string)
    if !ok || ownerID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    fileID := r.URL.Query().Get("id")
    if fileID == "" {
        http.Error(w, "missing file id", http.StatusBadRequest)
        return
    }

    shareToken, err := h.Service.CreateShare(models.ShareRequest{FileID: fileID}, ownerID)
    if err != nil {
        http.Error(w, "create share error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "link": "/public/file?token=" + shareToken,
    })
}

// GET /share?id=shareID
func (h *ShareHandler) AccessShare(w http.ResponseWriter, r *http.Request) {
    requesterID, _ := r.Context().Value("user_id").(string)
    shareID := r.URL.Query().Get("id")

    share, err := h.Service.AccessShare(shareID, requesterID)
    if err != nil {
        http.Error(w, "access denied: "+err.Error(), http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(share)
}

// GET /share/stats?file_id=123
func (h *ShareHandler) PublicStats(w http.ResponseWriter, r *http.Request) {
    fileID := r.URL.Query().Get("file_id")
    stats, err := h.Service.GetPublicStats(fileID)
    if err != nil {
        http.Error(w, "stats not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(stats)
}