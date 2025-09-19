package handlers

import (
    "encoding/json"
    "net/http"
    "balkanid-capstone/internal/models"
    "balkanid-capstone/internal/services"
)

type AuthHandler struct {
    Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
    return &AuthHandler{Service: service}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    userID, err := h.Service.Signup(creds)
    if err != nil {
        http.Error(w, "signup failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message": "signup successful",
        "user_id": userID,
    })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    token, role, err := h.Service.Login(creds)
    if err != nil {
        http.Error(w, "invalid email or password", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
        "role":  role,
    })
}