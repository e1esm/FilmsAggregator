package reindexer

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
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

func (cr *CacheRepository) Add(ctx context.Context, film dbModel.Film) (dbModel.Film, error) {
	film.CacheTime = time.Now()
	cr.db.WithContext(ctx)
	_, err := cr.db.Insert(cr.namespace, film, "id=serial()")
	if err != nil {
		logger.Logger.Error(err.Error(), zap.String("film", film.Title))
		return dbModel.Film{}, err
	}
	return film, nil
}

func (cr *CacheRepository) FindByName(ctx context.Context, name string) ([]*dbModel.Film, error) {
	films := make([]*dbModel.Film, 0)
	cr.db.WithContext(ctx)
	received, err := cr.db.Query(cr.namespace).Where("title", reindexer.EQ, name).Exec().FetchAll()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(received); i++ {
		films = append(films, received[i].(*dbModel.Film))
		if err != nil {
			logger.Logger.Error(err.Error())
			return nil, err
		}
	}

	return films, nil
}

func (cr *CacheRepository) Delete(ctx context.Context, request api.DeleteRequest) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE \"title\" = '%s' AND \"genre\"= '%s' AND \"released_year\"= %d;", cr.namespace, request.Title, request.Genre, request.ReleasedYear)
	cr.db.WithContext(ctx)
	_, err := cr.db.ExecSQL(query).FetchOne()
	if err != nil {
		return err
	}
	return nil
}

func (cr *CacheRepository) Verify(ctx context.Context, film *dbModel.Film) bool {
	alreadyExists := false
	cr.db.WithContext(ctx)
	_, err := cr.db.Query(cr.namespace).Where("hashcode", reindexer.EQ, film.HashCode).Exec().FetchOne()
	switch {
	case err == reindexer.ErrNotFound:
		return alreadyExists
	case err != nil:
		logger.Logger.Error(err.Error())
		return alreadyExists
	default:
		alreadyExists = true
		return alreadyExists
	}
}

func (cr *CacheRepository) DeleteCachedWithCtx(ctx context.Context) error {
	deleteRequest := ctx.Value("request")
	return cr.Delete(ctx, deleteRequest.(api.DeleteRequest))
}
