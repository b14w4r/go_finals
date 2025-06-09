package services

import (
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"

    "banking_service_project/models"
    "banking_service_project/repositories"
)

type AuthService interface {
    Register(username, email, password string) (*models.User, error)
    Login(email, password string) (string, error)
    ParseToken(tokenStr string) (string, error)
}

type authService struct {
    userRepo  repositories.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
    return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
    // Check uniqueness
    if _, err := s.userRepo.GetByEmail(email); err == nil {
        return nil, errors.New("email already in use")
    }
    if _, err := s.userRepo.GetByUsername(username); err == nil {
        return nil, errors.New("username already in use")
    }

    hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Username: username,
        Email:    email,
        Password: string(hashedPass),
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Subject:   string(rune(user.ID)),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    })

    tokenStr, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }
    return tokenStr, nil
}

func (s *authService) ParseToken(tokenStr string) (string, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    if err != nil || !token.Valid {
        return "", errors.New("invalid token")
    }
    claims := token.Claims.(*jwt.RegisteredClaims)
    return claims.Subject, nil
}
