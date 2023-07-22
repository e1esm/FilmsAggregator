package reindexer

import (
	"context"
	"fmt"
	dbModel "github.com/e1esm/FilmsAggregator/internal/models/db"
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
	err = db.OpenNamespace(config.Reindexer.DBName, reindexer.DefaultNamespaceOptions(), dbModel.Film{})
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}
	return &CacheRepository{db: db, namespace: config.Reindexer.DBName}
}

func (cr *CacheRepository) Add(ctx context.Context, film *dbModel.Film) (dbModel.Film, error) {
	film.CacheTime = time.Now()
	cr.db.WithContext(ctx)
	_, err := cr.db.Insert(cr.namespace, film, "id=serial()")
	if err != nil {
		logger.Logger.Error(err.Error(), zap.String("film", film.Title))
		return dbModel.Film{}, err
	}
	return *film, nil
}

func (cr *CacheRepository) FindByName(ctx context.Context, name string) ([]*dbModel.Film, error) {
	films := make([]*dbModel.Film, 0)
	received, err := cr.db.Query(cr.namespace).Where("title", reindexer.EQ, name).Exec().FetchAll()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(received); i++ {
		films = append(films, received[i].(*dbModel.Film))
		//err = json.Unmarshal([]byte(received[i].(dbModel.Film)), films[i])
		if err != nil {
			logger.Logger.Error(err.Error())
			return nil, err
		}
	}

	return films, nil
}

func (cr *CacheRepository) Delete(ctx context.Context, name string) (dbModel.Film, error) {
	return dbModel.Film{}, nil
}
