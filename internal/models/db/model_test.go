package db

import (
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding(t *testing.T) {
	film := &Film{Title: "X", Crew: &general.Crew{}, Revenue: 100, ReleasedYear: 2004}
	encode(film)
	filmReceived := NewFilm(uuid.UUID{}, film.Title, film.Crew, film.ReleasedYear, film.Revenue, film.Genre)

	assert.Equal(t, film.HashCode, filmReceived.HashCode)
}
