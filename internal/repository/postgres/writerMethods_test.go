package postgres

import (
	"context"
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
				uuid.New(), "XXX", &general.Crew{Producers: []*general.Producer{{Person: general.NewPerson(uuid.New(), "first", "2000-06-01", "m")}}, Actors: []*general.Actor{{Person: general.NewPerson(uuid.New(), "first", "2004-03-21", "m"), Role: "Kicker"}}}, 2004, 199.99, "fantasy",
			),
			status: SUCCESS,
		},
	}

	for _, test := range testTable {
		received, _ := testRepository.Add(context.Background(), test.film)
		assert.Equal(t, test.film, received)

	}

}
