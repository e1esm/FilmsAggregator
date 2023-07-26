package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/auth"
	"github.com/e1esm/FilmsAggregator/internal/repository/authentication"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

const (
	salt       = "rknfgnrgnrnkgnrgnrgn"
	signingKey = "grglkrgnrklr;gq;eremr'gmrgr"
	tokenTTl   = 12 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}

type AuthorizationService interface {
	CreateUser(ctx context.Context, user auth.User) (auth.User, error)
	GenerateToken(ctx context.Context, username string, password string) (string, error)
}

type AuthService struct {
	AuthRepository authentication.Authenticator
}

func NewAuthService(repository authentication.Authenticator) AuthorizationService {
	return &AuthService{AuthRepository: repository}
}

func (as *AuthService) CreateUser(ctx context.Context, user auth.User) (auth.User, error) {

	user.Password = generateHashForThePassword(user.Password)
	user.ID = uuid.New()
	received, err := as.AuthRepository.CreateUser(ctx, user)
	if err != nil {
		return received, err
	}
	return received, nil
}

func (as *AuthService) GenerateToken(ctx context.Context, username string, password string) (string, error) {
	password = generateHashForThePassword(password)
	user, err := as.AuthRepository.GetUser(ctx, username, password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTl).Unix(),
		IssuedAt:  time.Now().Unix(),
	},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func generateHashForThePassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
