package service

import (
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository/authentication"
)

type AuthorizationService interface {
	CreateUser(user models.User) error
	GenerateToken(username string, password string) error
}

type AuthService struct {
	AuthRepository authentication.Authenticator
}

func (a AuthService) CreateUser(user models.User) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthService) GenerateToken(username string, password string) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(repository authentication.Authenticator) AuthorizationService {
	return AuthService{AuthRepository: repository}
}
