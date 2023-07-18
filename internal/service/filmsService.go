package service

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository"
)

type Service interface {
	AddFilm(context.Context, models.Film) (models.Film, error)
}

type FilmsService struct {
	Repositories *repository.Repositories
}

func NewService(repositories *repository.Repositories) *FilmsService {
	return &FilmsService{Repositories: repositories}
}
