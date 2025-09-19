package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

type AuthHandler struct {
    DB *sql.DB
}

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// Signup handler
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "hash error", http.StatusInternalServerError)
        return
    }

    var userID string
    err = h.DB.QueryRow(`
        INSERT INTO users (email, password_hash) 
        VALUES ($1, $2) RETURNING id
    `, creds.Email, string(hashed)).Scan(&userID)
    if err != nil {
        http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message": "signup successful",
        "user_id": userID,
    })
}

// Login handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    var userID string
    var passwordHash string
    err := h.DB.QueryRow(`
        SELECT id, password_hash FROM users WHERE email=$1
    `, creds.Email).Scan(&userID, &passwordHash)
    if err != nil {
        http.Error(w, "invalid email or password", http.StatusUnauthorized)
        return
    }

    if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)) != nil {
        http.Error(w, "invalid email or password", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "token error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}
