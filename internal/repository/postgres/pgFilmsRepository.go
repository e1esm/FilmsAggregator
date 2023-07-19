package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	uuidHash "github.com/google/uuid"
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

func (fr *FilmsRepository) Add(ctx context.Context, film *models.Film) (models.Film, error) {
	uuid.GenerateUUIDs(film)

	tx, err := fr.Pool.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Couldn't have begun transaction", zap.String("err", err.Error()))
	}
	fr.TransactionManager.Add(film.ID, tx)
	defer func() {
		tx.Rollback(ctx)
		fr.TransactionManager.Delete(film.ID)
	}()
	if err = fr.AddFilm(ctx, film); err != nil {
		return models.Film{}, err
	}

	if err = fr.AddCrew(ctx, film); err != nil {
		return models.Film{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Film{}, err
	}
	return *film, nil
}

func (fr *FilmsRepository) AddFilm(ctx context.Context, film *models.Film) error {
	tx := fr.TransactionManager.Get(film.ID)
	if tx == nil {
		fr.TransactionManager.Delete(film.ID)
		return fmt.Errorf("tracsaction wasn't started, neither was deleted")
	}
	_, err := tx.Exec(ctx, "INSERT INTO film (id, title, release_year, revenue) VALUES ($1, $2, $3, $4);",
		film.ID,
		film.Title,
		film.ReleasedYear,
		film.Revenue)
	if err != nil {
		logger.Logger.Error("Couldn't have inserted film",
			zap.String("title", film.Title),
			zap.String("err", err.Error()))
		return err
	}
	return nil
}

//AddCrew

func (fr *FilmsRepository) AddCrew(ctx context.Context, film *models.Film) error {

	if err := fr.AddProducers(ctx, film.ID, film.Crew.Producers); err != nil {
		return err
	}

	if err := fr.AddActors(ctx, film.ID, film.Crew.Actors); err != nil {
		return err
	}
	return nil
}

func (fr *FilmsRepository) AddProducers(ctx context.Context, id uuidHash.UUID, producers []models.Producer) error {
	tx := fr.TransactionManager.Get(id)
	if tx == nil {
		fr.TransactionManager.Delete(id)
		return fmt.Errorf("tracsaction wasn't started, neither was deleted")
	}
	for _, producer := range producers {
		_, err := tx.Exec(ctx, "INSERT INTO producer (id, name, birthdate, gender) VALUES ($1, $2, $3, $4);",
			producer.ID,
			producer.Name,
			producer.Birthdate,
			producer.Gender)
		if err != nil {
			logger.Logger.Error("Couldn't have inserted producer",
				zap.String("producer", producer.Name),
				zap.String("err", err.Error()))
			return err
		}
	}
	return nil
}

func (fr *FilmsRepository) AddActors(ctx context.Context, id uuidHash.UUID, actors []models.Actor) error {
	tx := fr.TransactionManager.Get(id)
	if tx == nil {
		fr.TransactionManager.Delete(id)
		return fmt.Errorf("tracsaction wasn't started, neither was deleted")
	}
	for _, actor := range actors {
		_, err := tx.Exec(ctx, "INSERT INTO actor (id, name, birthdate, gender, role) VALUES ($1, $2, $3, $4, $5);",
			actor.ID,
			actor.Name,
			actor.Birthdate,
			actor.Gender,
			actor.Role)
		if err != nil {
			logger.Logger.Error("Couldn't have inserted actor",
				zap.String("actor", actor.Name),
				zap.String("err", err.Error()))
			return err
		}
	}
	return nil
}
