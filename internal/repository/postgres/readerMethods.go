package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgtype"
)

var (
	queryError    = errors.New("error occurred while querying from DB")
	scanningError = errors.New("error occurred while scanning values from DB")
)

func (fr *FilmsRepository) Verify(ctx context.Context, film *db.Film) bool {
	doesAlreadyExist := false
	var title string
	query := "select title from film where hashcode = $1;"
	err := fr.Pool.QueryRow(ctx, query, film.HashCode).Scan(&title)
	switch {
	case err == sql.ErrNoRows:
		return doesAlreadyExist
	case err != nil:
		logger.Logger.Error(err.Error())
	default:
		doesAlreadyExist = true
		return doesAlreadyExist
	}

	return doesAlreadyExist
}

func (fr *FilmsRepository) FindByName(ctx context.Context, name string) ([]*db.Film, error) {
	foundFilms := make([]*db.Film, 0)
	query := "select * from film where title = $1;"
	rows, err := fr.Pool.Query(ctx, query, name)
	if err != nil {
		return nil, queryError
	}
	i := -1
	for rows.Next() {
		i++
		foundFilms = append(foundFilms, &db.Film{})
		if err = rows.Scan(&foundFilms[i].ID, &foundFilms[i].Title, &foundFilms[i].Genre, &foundFilms[i].ReleasedYear, &foundFilms[i].Revenue); err != nil {
			return nil, scanningError
		}
		foundFilms[i].Crew.Actors = make([]*general.Actor, 0)
		foundFilms[i].Crew.Producers = make([]*general.Producer, 0)
	}

	films, err := fr.findCrew(ctx, foundFilms)
	if err != nil {
		return nil, err
	}
	return films, nil
}

func (fr *FilmsRepository) findCrew(ctx context.Context, films []*db.Film) ([]*db.Film, error) {
	query := `WITH Producers AS (
		SELECT
	producer.id AS producer_id,
		producer.name AS producer_name,
		producer.gender AS producer_gender,
		CAST(producer.birthdate AS TEXT) AS producer_birthdate
	FROM
	producer
	INNER JOIN crew c ON c.producer_id = producer.id
	WHERE
	c.film_id = $1::uuid
	),
	Actors AS (
		SELECT
	actor.id AS actor_id,
		actor.name AS actor_name,
		actor.gender AS actor_gender,
		cast(actor.birthdate AS TEXT) AS actor_birthdate,
		actor.role AS actor_role
	FROM
	actor
	INNER JOIN crew c ON c.actor_id = actor.id
	WHERE
	c.film_id = $2::uuid
	)
	SELECT
	*
		FROM
	Producers
	FULL OUTER JOIN Actors ON Producers.producer_id = Actors.actor_id;`
	tx, err := fr.Pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	for _, film := range films {
		rows, err := tx.Query(ctx, query, film.ID, film.ID)
		if err != nil {
			return nil, err
		}
		tempResponse := &db.ResponseTemp{}
		for rows.Next() {
			err = rows.Scan(&tempResponse.ProducerID,
				&tempResponse.ProducerName,
				&tempResponse.ProducerGender,
				&tempResponse.ProducerBirthdate,
				&tempResponse.ActorID,
				&tempResponse.ActorName,
				&tempResponse.ActorGender,
				&tempResponse.ActorBirthdate,
				&tempResponse.ActorRole)
			if err != nil {
				logger.Logger.Error(err.Error())
				return nil, err
			}
			if tempResponse.ActorID.Status == pgtype.Null {
				film.Crew.Producers = append(film.Crew.Producers,
					general.NewProducer(tempResponse.ProducerID.Bytes,
						tempResponse.ProducerName.String,
						tempResponse.ProducerBirthdate.String,
						tempResponse.ProducerGender.String))
			} else {
				film.Crew.Actors = append(film.Crew.Actors,
					general.NewActor(tempResponse.ActorID.Bytes,
						tempResponse.ActorName.String,
						tempResponse.ActorBirthdate.String,
						tempResponse.ActorGender.String,
						tempResponse.ActorRole.String))
			}
		}
	}
	return films, nil
}
