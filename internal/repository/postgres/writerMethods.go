package postgres

import (
	"context"
	"errors"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	uuidHash "github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	transactionError = errors.New("transaction wasn't started, neither it'd been deleted")
)

func (fr *FilmsRepository) Add(ctx context.Context, film db.Film) (db.Film, error) {
	tx, err := fr.Pool.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Couldn't have begun transaction", zap.String("err", err.Error()))
	}
	fr.TransactionManager.Add(film.ID, tx)
	defer func() {
		tx.Rollback(ctx)
		fr.TransactionManager.Delete(film.ID)
	}()
	if err = fr.addFilm(ctx, &film); err != nil {
		return db.Film{}, err
	}

	if err = fr.addWorkers(ctx, &film); err != nil {
		return db.Film{}, err
	}
	if err = fr.addCrew(ctx, &film); err != nil {
		return db.Film{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return db.Film{}, err
	}
	return film, nil
}

func (fr *FilmsRepository) addFilm(ctx context.Context, film *db.Film) error {
	tx := fr.TransactionManager.Get(film.ID)
	if tx == nil {
		fr.TransactionManager.Delete(film.ID)
		return transactionError
	}
	_, err := tx.Exec(ctx, "INSERT INTO film (id, title, release_year, revenue, hashcode) VALUES ($1, $2, $3, $4, $5);",
		film.ID,
		film.Title,
		film.ReleasedYear,
		film.Revenue,
		film.HashCode)
	if err != nil {
		logger.Logger.Error("Couldn't have inserted film",
			zap.String("title", film.Title),
			zap.String("err", err.Error()))
		return err
	}
	return nil
}

func (fr *FilmsRepository) addWorkers(ctx context.Context, film *db.Film) error {

	if err := fr.addProducers(ctx, film.ID, film.Crew.Producers); err != nil {
		return err
	}

	if err := fr.addActors(ctx, film.ID, film.Crew.Actors); err != nil {
		return err
	}
	return nil
}

func (fr *FilmsRepository) addProducers(ctx context.Context, id uuidHash.UUID, producers []*general.Producer) error {
	tx, err := fr.TransactionManager.VerifyAndGet(id)
	if err != nil {
		return err
	}
	for _, producer := range producers {
		_, err = tx.Exec(ctx, "INSERT INTO producer (id, name, birthdate, gender) VALUES ($1, $2, $3, $4);",
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

func (fr *FilmsRepository) addActors(ctx context.Context, id uuidHash.UUID, actors []*general.Actor) error {
	tx, err := fr.TransactionManager.VerifyAndGet(id)
	if err != nil {
		return err
	}
	for _, actor := range actors {
		_, err = tx.Exec(ctx, "INSERT INTO actor (id, name, birthdate, gender, role) VALUES ($1, $2, $3, $4, $5);",
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

func (fr *FilmsRepository) addCrew(ctx context.Context, film *db.Film) error {
	tx, err := fr.TransactionManager.VerifyAndGet(film.ID)
	if err != nil {
		return err
	}
	for _, producer := range film.Crew.Producers {
		_, err := tx.Exec(ctx, "INSERT INTO crew (producer_id, film_id) VALUES ($1, $2);", producer.ID, film.ID)
		if err != nil {
			return err
		}
	}

	for _, actor := range film.Crew.Actors {
		_, err = tx.Exec(ctx, "INSERT INTO crew (actor_id, film_id) VALUES ($1, $2)", actor.ID, film.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
