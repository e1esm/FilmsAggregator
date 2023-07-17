package postgres

import (
	"database/sql"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FilmsRepository struct {
	DB *sql.DB
}

func NewFilmsRepository(cfg config.Config) *FilmsRepository {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.ContainerName,
		cfg.Postgres.Port,
		cfg.Postgres.DatabaseName)
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		logger.Logger.Fatal("Couldn't have opened connection with DB", zap.String("err", err.Error()))
	}
	return &FilmsRepository{DB: db}
}

func (fr *FilmsRepository) Add(filmd models.Film) uuid.UUID {
	return uuid.New()
}
