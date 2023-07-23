package db

import (
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding(t *testing.T) {
	film := &Film{ID: uuid.New(), Title: "X", Crew: general.Crew{}, Revenue: 100, ReleasedYear: 2004}
	encode(film)
	filmReceived := NewFilm(film.ID, film.Title, film.Crew, film.ReleasedYear, film.Revenue)

	assert.Equal(t, film.HashCode, filmReceived.HashCode)
}
