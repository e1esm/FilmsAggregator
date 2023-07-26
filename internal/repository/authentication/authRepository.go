package authentication

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Authenticator interface {
}

type AuthRepository struct {
	Pool *pgxpool.Pool
}

func NewAuthRepository(config config.Config) Authenticator {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?pool_max_conns=%d",
		config.AuthDB.User,
		config.AuthDB.Password,
		config.AuthDB.ContainerName,
		config.AuthDB.Port,
		config.AuthDB.DatabaseName,
		config.AuthDB.Connections)
	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		logger.Logger.Fatal("Couldn't have opened connection with DB", zap.String("err", err.Error()))
	}
	return AuthRepository{Pool: pool}
}
