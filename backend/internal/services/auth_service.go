package services

import (
	"balkanid-capstone/internal/models"
	"balkanid-capstone/internal/repo"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

type AuthService struct {
    Repo *repo.UserRepo
}

func NewAuthService(r *repo.UserRepo) *AuthService {
    return &AuthService{Repo: r}
}

func (s *AuthService) Signup(creds models.Credentials) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return s.Repo.CreateUser(creds.Name, creds.Email, string(hashed))
}

func (s *AuthService) Login(creds models.Credentials) (string, string, error) {
    user, err := s.Repo.GetUserByEmail(creds.Email)
    if err != nil {
        return "", "", err
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) != nil {
        return "", "", sql.ErrNoRows
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &models.Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", "", err
    }

    return tokenString, user.Role, nil
}