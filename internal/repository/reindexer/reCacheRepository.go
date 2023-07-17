package reindexer

import (
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/google/uuid"
)

type CacheRepository struct {
}

func NewFilmsRepository() *CacheRepository {
	return nil
}

func (cr *CacheRepository) Add(film models.Film) uuid.UUID {
	return uuid.New()
}
