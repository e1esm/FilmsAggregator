package authentication

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type Authenticator interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
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
	return &AuthRepository{Pool: pool}
}

func (ar *AuthRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := ar.Pool.Exec(ctx, "INSERT INTO users (id, username, password, role) VALUES ($1, $2, $3, $4);",
		user.ID, user.Username, user.Password, user.Role)
	if err != nil {
		logger.Logger.Error(err.Error())
		return models.User{}, err
	}
	return user, nil
}

func (ar *AuthRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var user models.User
	row := ar.Pool.QueryRow(ctx, "SELECT * FROM users WHERE username = $1 AND password = $2", username, password)
	if err := row.Scan(&user); err != nil {
		logger.Logger.Error(err.Error())
		return models.User{}, err
	}
	return user, nil
}
