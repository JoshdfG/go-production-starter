package usecase

import (
	"fmt"
	"time"

	"todo-clean/internal/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo      UserRepository
	jwtSecret string
	expiry    time.Duration
}

func NewAuthUseCase(repo UserRepository, jwtSecret string, expiryHours int) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
		expiry:    time.Duration(expiryHours) * time.Hour,
	}
}

func (uc *AuthUseCase) Register(email, password string) (*entity.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	user := &entity.User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	return user, nil
}

func (uc *AuthUseCase) Login(email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(uc.expiry).Unix(),
		"iat": time.Now().Unix(),
	})

	signed, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("Login: %w", err)
	}

	return signed, nil
}

func (uc *AuthUseCase) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(uc.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid subject")
	}

	return userID, nil
}
