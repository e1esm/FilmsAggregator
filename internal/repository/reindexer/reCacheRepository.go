package reindexer

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/restream/reindexer/v3"
)

type CacheRepository struct {
	DB *reindexer.Reindexer
}

func NewCacheRepository(config config.Config) *CacheRepository {
	dsn := fmt.Sprintf("cproto://%s:%d/%s", config.Reindexer.ContainerName, config.Reindexer.Port, config.Reindexer.DBName)
	db := reindexer.NewReindex(dsn)
	return &CacheRepository{DB: db}
}

func (cr *CacheRepository) Add(ctx context.Context, film *models.Film) (models.Film, error) {
	return models.Film{}, nil
}

func (cr *CacheRepository) GetByName(ctx context.Context, name string) ([]*models.Film, error) {
	return nil, nil
}
