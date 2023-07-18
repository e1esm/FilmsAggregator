package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type FilmsRepository struct {
	Pool *pgxpool.Pool
}

func NewFilmsRepository(cfg config.Config) repository.Repository {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?pool_max_conns=%d",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.ContainerName,
		cfg.Postgres.Port,
		cfg.Postgres.DatabaseName,
		cfg.Postgres.Connections)
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Logger.Fatal("Couldn't have opened connection with DB", zap.String("err", err.Error()))
	}
	return &FilmsRepository{Pool: pool}
}

func (fr *FilmsRepository) Add(ctx context.Context, film *models.Film) (models.Film, error) {
	uuid.GenerateUUIDs(film)

	tx, err := fr.Pool.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Couldn't have begun transaction", zap.String("err", err.Error()))
	}
	defer tx.Rollback(ctx)
	if err != nil {
		logger.Logger.Error("Couldn't have started transaction", zap.String("err", err.Error()))
		return models.Film{}, nil
	}
	_, err = tx.Exec(ctx, "INSERT INTO film (id, title, release_year, revenue) VALUES ($1, $2, $3, $4);",
		film.ID,
		film.Title,
		film.ReleasedYear,
		film.Revenue)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Film{}, err
	}
	return *film, nil
}
