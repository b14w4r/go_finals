package repositories

import (
    "database/sql"
    "errors"
    "time"

    "banking_service_project/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
    GetByUsername(username string) (*models.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    query := `INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
    user.CreatedAt = time.Now()
    err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, username, email, password, created_at FROM users WHERE email=$1`
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, username, email, password, created_at FROM users WHERE username=$1`
    err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}
