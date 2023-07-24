package uuid

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/google/uuid"
)

func GenerateUUIDs(film *db.Film) {
	for i := 0; i < len(film.Crew.Producers); i++ {
		film.Crew.Producers[i].ID = uuid.New()
	}
	for i := 0; i < len(film.Crew.Actors); i++ {
		film.Crew.Actors[i].ID = uuid.New()
	}
}
