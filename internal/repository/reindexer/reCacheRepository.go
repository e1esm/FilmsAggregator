package reindexer

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"go.uber.org/zap"
	"time"
)

var (
	retires = 10
)

type CacheRepository struct {
	namespace string
	db        *reindexer.Reindexer
}

func NewCacheRepository(config config.Config) *CacheRepository {
	dsn := fmt.Sprintf("cproto://%s:%d/%s", config.Reindexer.ContainerName, config.Reindexer.Port, config.Reindexer.DBName)
	db := reindexer.NewReindex(dsn, reindexer.WithCreateDBIfMissing())
	var err error
	try := 0
	ticker := time.NewTicker(1 * time.Second)
	for try < retires {
		select {
		case <-ticker.C:
			err = db.Ping()
			if err == nil {
				break
			}
			db = reindexer.NewReindex(dsn, reindexer.WithCreateDBIfMissing())
		}
		try++
	}
	if err != nil {
		return nil
	}
	err = db.OpenNamespace(config.Reindexer.DBName, reindexer.DefaultNamespaceOptions(), api.Film{})
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}
	return &CacheRepository{db: db, namespace: config.Reindexer.DBName}
}

func (cr *CacheRepository) Add(ctx context.Context, film *api.Film) (api.Film, error) {
	cr.db.WithContext(ctx)
	_, err := cr.db.Insert(cr.namespace, film, "id=serial()")
	if err != nil {
		logger.Logger.Error(err.Error(), zap.String("film", film.Title))
		return api.Film{}, err
	}
	return api.Film{}, nil
}

func (cr *CacheRepository) FindByName(ctx context.Context, name string) ([]*api.Film, error) {
	return nil, nil
}

func (cr *CacheRepository) Delete(ctx context.Context, name string) (api.Film, error) {
	return api.Film{}, nil
}
