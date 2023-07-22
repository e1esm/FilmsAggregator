package api

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
)

type Film struct {
	ID           uuid.UUID    `json:"-" reindex:"-"`
	Title        string       `json:"title" reindex:"title,tree"`
	Crew         general.Crew `json:"crew"`
	ReleasedYear int          `json:"released_year" reindex:"released_year"`
	Revenue      float64      `json:"revenue" reindex:"revenue"`
}

func NewFilm(film db.Film) *Film {
	return &Film{ID: film.ID,
		Title:        film.Title,
		Crew:         film.Crew,
		Revenue:      film.Revenue,
		ReleasedYear: film.ReleasedYear,
	}
}
