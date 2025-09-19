package models

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

type User struct {
    ID           string
    Name         string
    Email        string
    PasswordHash string
    Role         string
}

type File struct {
    ID        string `json:"id"`
    Filename  string `json:"filename"`
    MimeType  string `json:"mime_type"`
    Size      int64  `json:"size"`
    CreatedAt string `json:"created_at"`
    Hash      string `json:"hash"`
    RefCount  int    `json:"ref_count"`
}

type FileBlob struct {
    ID          string
    Hash        string
    RefCount    int
    StoragePath string
}

type UploadResult struct {
    Hash     string `json:"hash"`
    Filename string `json:"filename"`
}