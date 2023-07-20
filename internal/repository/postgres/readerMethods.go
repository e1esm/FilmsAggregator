package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgtype"
)

func (fr *FilmsRepository) FindByName(ctx context.Context, name string) ([]*models.Film, error) {
	foundFilms := make([]*models.Film, 0)
	query := "select * from film where title = $1;"
	rows, err := fr.Pool.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("error occurred while querying from DB")
	}
	i := -1
	for rows.Next() {
		i++
		foundFilms = append(foundFilms, &models.Film{})
		if err = rows.Scan(&foundFilms[i].ID, &foundFilms[i].Title, &foundFilms[i].ReleasedYear, &foundFilms[i].Revenue); err != nil {
			return nil, fmt.Errorf("error occurred while scanning fetched values")
		}
		foundFilms[i].Crew.Actors = make([]models.Actor, 0)
		foundFilms[i].Crew.Producers = make([]models.Producer, 0)
	}

	return fr.findCrew(ctx, foundFilms)
}

func (fr *FilmsRepository) findCrew(ctx context.Context, films []*models.Film) ([]*models.Film, error) {
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
		var (
			producerID        pgtype.UUID
			producerName      pgtype.Text
			producerGender    pgtype.Text
			producerBirthdate pgtype.Text
			actorID           pgtype.UUID
			actorName         pgtype.Text
			actorGender       pgtype.Text
			actorBirthdate    pgtype.Text
			actorRole         pgtype.Text
		)
		for rows.Next() {
			err = rows.Scan(&producerID, &producerName, &producerGender, &producerBirthdate,
				&actorID, &actorName, &actorGender, &actorBirthdate, &actorRole)
			if err != nil {
				logger.Logger.Error(err.Error())
				return nil, err
			}
			if actorID.Status == pgtype.Null {
				film.Crew.Producers = append(film.Crew.Producers, *models.NewProducer(producerID.Bytes, producerName.String, producerBirthdate.String, producerGender.String))
			} else {
				film.Crew.Actors = append(film.Crew.Actors, *models.NewActor(actorID.Bytes, actorName.String, actorBirthdate.String, actorGender.String, actorRole.String))
			}
		}
	}

	return films, nil
}
