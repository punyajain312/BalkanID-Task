package repo

import (
    "database/sql"
    "balkanid-capstone/internal/models"
)

type UserRepo struct {
    DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
    return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(name, email, passwordHash string) (string, error) {
    var userID string
    err := r.DB.QueryRow(`
        INSERT INTO users (name, email, password_hash, role)
        VALUES ($1, $2, $3, 'user')
        RETURNING id
    `, name, email, passwordHash).Scan(&userID)
    return userID, err
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.DB.QueryRow(`
        SELECT id, name, email, password_hash, role 
        FROM users WHERE email=$1
    `, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

    if err != nil {
        return nil, err
    }
    return &user, nil
}