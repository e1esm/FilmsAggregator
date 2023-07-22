package service

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"go.uber.org/zap"
	"time"
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

	received, err := fs.Repositories.CacheRepo.FindByName(ctx, name)
	if err != nil {
		logger.Logger.Error("Couldn't have retrieved films from cache", zap.String("err", err.Error()))
	}
	current := time.Now()
	// TODO Change Hardcoded value
	isUpToDate := true
	for i := 0; i < len(received); i++ {
		if current.Sub(received[i].CacheTime).Minutes() > 15 {
			isUpToDate = false
		}
	}

	if isUpToDate {
		return received, nil
	} else {
		for i := 0; i < len(received); i++ {
			_, err = fs.Repositories.CacheRepo.Delete(ctx, received[i].Title)
			if err != nil {
				logger.Logger.Error("Couldn't have deleted film from cache",
					zap.String("err", err.Error()),
					zap.String("film", received[i].Title))
			}
		}
	}
	received, err = fs.Repositories.MainRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return received, nil
}
