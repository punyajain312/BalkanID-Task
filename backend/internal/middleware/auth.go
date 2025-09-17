package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("supersecretkey")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "missing auth header", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return JwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        // Safely extract values from claims
        userId, ok := claims["user_id"].(string)
        if !ok {
            http.Error(w, "invalid token claims", http.StatusUnauthorized)
            return
        }

        role, _ := claims["role"].(string) // default to "" if not present

        // Pass values into request context
        ctx := context.WithValue(r.Context(), "user_id", userId)
        ctx = context.WithValue(ctx, "role", role)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
