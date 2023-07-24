package api

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
)

type Film struct {
	Title        string        `json:"title"`
	Crew         *general.Crew `json:"crew"`
	ReleasedYear int           `json:"released_year"`
	Revenue      float64       `json:"revenue"`
	Genre        string
}

func NewFilm(film db.Film) *Film {
	return &Film{
		Title:        film.Title,
		Crew:         film.Crew,
		Revenue:      film.Revenue,
		ReleasedYear: film.ReleasedYear,
		Genre:        film.Genre,
	}
}
