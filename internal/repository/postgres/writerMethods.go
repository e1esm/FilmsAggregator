package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	uuidHash "github.com/google/uuid"
	"go.uber.org/zap"
)

func (fr *FilmsRepository) Add(ctx context.Context, film *api.Film) (api.Film, error) {
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
	if err = fr.addFilm(ctx, film); err != nil {
		return api.Film{}, err
	}

	if err = fr.addWorkers(ctx, film); err != nil {
		return api.Film{}, err
	}
	if err = fr.addCrew(ctx, film); err != nil {
		return api.Film{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return api.Film{}, err
	}
	return *film, nil
}

func (fr *FilmsRepository) addFilm(ctx context.Context, film *api.Film) error {
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

func (fr *FilmsRepository) addWorkers(ctx context.Context, film *api.Film) error {

	if err := fr.addProducers(ctx, film.ID, film.Crew.Producers); err != nil {
		return err
	}

	if err := fr.addActors(ctx, film.ID, film.Crew.Actors); err != nil {
		return err
	}
	return nil
}

func (fr *FilmsRepository) addProducers(ctx context.Context, id uuidHash.UUID, producers []api.Producer) error {
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

func (fr *FilmsRepository) addActors(ctx context.Context, id uuidHash.UUID, actors []api.Actor) error {
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

func (fr *FilmsRepository) addCrew(ctx context.Context, film *api.Film) error {
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
