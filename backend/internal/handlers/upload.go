package handlers

import (
    "encoding/json"
    "net/http"

    "balkanid-capstone/internal/models"
    "balkanid-capstone/internal/services"
)

type UploadHandler struct {
    Service *services.UploadService
}

func NewUploadHandler(service *services.UploadService) *UploadHandler {
    return &UploadHandler{Service: service}
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form (max 20MB)
    if err := r.ParseMultipartForm(20 << 20); err != nil {
        http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
        return
    }

    files := r.MultipartForm.File["files"]
    if len(files) == 0 {
        http.Error(w, "no files uploaded", http.StatusBadRequest)
        return
    }

    userID, ok := r.Context().Value("user_id").(string)
    if !ok || userID == "" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var uploaded []models.UploadResult

    for _, header := range files {
        f, err := header.Open()
        if err != nil {
            http.Error(w, "failed to open file", http.StatusInternalServerError)
            return
        }
        defer f.Close()

        result, err := h.Service.UploadFile(
            userID,
            f,
            header.Filename,
            header.Header.Get("Content-Type"),
            header.Size,
        )
        if err != nil {
            http.Error(w, "upload error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        uploaded = append(uploaded, result)
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "upload successful",
        "files":   uploaded,
    })
}