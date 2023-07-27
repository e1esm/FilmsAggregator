package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/auth"
	"github.com/e1esm/FilmsAggregator/internal/repository/authentication"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	salt       string
	signingKey string
	TokenTTL   time.Duration
)

func init() {
	err := godotenv.Load("auth.env")
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	TokenTTL, err = time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		TokenTTL = 3600000000000
	}
	salt = os.Getenv("SALT")
	signingKey = os.Getenv("SIGNING_KEY")
	logger.Logger.Info("Auth Variables", zap.String("salt", salt),
		zap.String("signingKey", signingKey),
		zap.String("tokenTTL", fmt.Sprintf("%v", TokenTTL)))
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID   uuid.UUID `json:"user_id"`
	UserRole string    `json:"user_role"`
}

type AuthorizationService interface {
	CreateUser(ctx context.Context, user auth.User) (auth.User, error)
	GenerateToken(ctx context.Context, username string, password string) (string, error)
	ParseToken(ctx context.Context, token string) (uuid.UUID, string, error)
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
		ExpiresAt: time.Now().Add(TokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	},
		user.ID,
		string(user.Role),
	})

	return token.SignedString([]byte(signingKey))
}

func (as *AuthService) ParseToken(ctx context.Context, token string) (uuid.UUID, string, error) {
	receivedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return uuid.UUID{}, "", err
	}
	claims, ok := receivedToken.Claims.(*TokenClaims)
	if !ok {
		return uuid.UUID{}, "", errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserID, claims.UserRole, nil
}

func generateHashForThePassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
