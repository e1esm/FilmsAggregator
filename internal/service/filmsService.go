package service

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository"
)

type Service interface {
	Add(context.Context, *models.Film) (models.Film, error)
	Get(ctx context.Context, name string) ([]*models.Film, error)
}

type FilmsService struct {
	Repositories *repository.Repositories
}

func NewService(repositories *repository.Repositories) *FilmsService {
	return &FilmsService{Repositories: repositories}
}

func (fs *FilmsService) Add(ctx context.Context, film *models.Film) (models.Film, error) {
	inserted, err := fs.Repositories.MainRepo.Add(ctx, film)
	if err != nil {

		return models.Film{}, err
	}
	return inserted, nil
}

func (fs *FilmsService) Get(ctx context.Context, name string) ([]*models.Film, error) {
	received, err := fs.Repositories.MainRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return received, nil
}
