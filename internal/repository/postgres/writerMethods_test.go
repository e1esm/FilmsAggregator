package postgres

import (
	"context"
	"database/sql"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilmsRepository_Add(t *testing.T) {

	testTable := []struct {
		film   db.Film
		status Status
	}{
		{
			film: *db.NewFilm(
				uuid.New(), "XXX", &general.Crew{Producers: []*general.Producer{{Person: general.NewPerson(uuid.New(), "first", "2000-06-01", "m")}}, Actors: []*general.Actor{{Person: general.NewPerson(uuid.New(), "first", "2004-03-21", "m"), Role: "Kicker"}}}, 2004, 100000, "fantasy",
			),
			status: SUCCESS,
		},
		{
			film: *db.NewFilm(
				uuid.New(), "YYY", &general.Crew{Producers: []*general.Producer{{Person: general.NewPerson(uuid.New(), "first", "2000-06-01", "m")}}, Actors: []*general.Actor{{Person: general.NewPerson(uuid.New(), "first", "2004-03-21", "m"), Role: "Kicker"}}}, 1999, 1000000, "drama",
			),
			status: SUCCESS,
		},
	}

	for _, test := range testTable {
		received, _ := testRepository.Add(context.Background(), test.film)
		assert.Equal(t, test.film, received)

	}

}

func TestFilmsRepository_Delete(t *testing.T) {
	testTable := []struct {
		deleteRequest api.DeleteRequest
		err           error
		status        Status
	}{
		{
			deleteRequest: api.DeleteRequest{Title: "YYY", Genre: "drama", ReleasedYear: 1999},
			status:        SUCCESS,
		},
		{
			deleteRequest: api.DeleteRequest{Title: "AAA", Genre: "horror", ReleasedYear: 1987},
			err:           sql.ErrNoRows,
			status:        FAIL,
		},
	}

	for _, test := range testTable {
		err := testRepository.Delete(context.Background(), test.deleteRequest)
		assert.Equal(t, test.err, err)
	}
}
