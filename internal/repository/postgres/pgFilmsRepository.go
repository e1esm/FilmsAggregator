package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type FilmsRepository struct {
	Pool               *pgxpool.Pool
	TransactionManager *TransactionManager
}

func NewFilmsRepository(cfg config.Config, manager *TransactionManager) repository.Repository {
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
	return &FilmsRepository{Pool: pool, TransactionManager: manager}
}
