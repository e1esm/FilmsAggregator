package uuid

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/google/uuid"
)

type Generator interface {
	Generate() uuid.UUID
	GenerateUUIDs(film db.Film) *db.Film
}

type UUIDGenerator struct {
}

func (gen *UUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}

func (gen *UUIDGenerator) GenerateUUIDs(film db.Film) *db.Film {
	if film.Crew == nil || film.Crew.Producers == nil || film.Crew.Actors == nil {
		return &film
	}
	for i := 0; i < len(film.Crew.Producers); i++ {
		film.Crew.Producers[i].ID = gen.Generate()
	}
	for i := 0; i < len(film.Crew.Actors); i++ {
		film.Crew.Actors[i].ID = gen.Generate()
	}
	return &film
}
