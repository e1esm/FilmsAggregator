package mocks

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/google/uuid"
)

type MockIDGenerator struct {
}

func (m *MockIDGenerator) Generate() uuid.UUID {
	return uuid.UUID{}
}
func (m *MockIDGenerator) GenerateUUIDs(film db.Film) *db.Film {
	return &film
}
