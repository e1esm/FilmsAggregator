package reindexer

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models"
)

type CacheRepository struct {
}

func NewFilmsRepository() *CacheRepository {
	return nil
}

func (cr *CacheRepository) Add(ctx context.Context, film *models.Film) (models.Film, error) {
	return models.Film{}, nil
}
