package api

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
)

// Film model info
// @Description model's being operated on
type Film struct {
	Title        string        `json:"title"`         //Title of the show
	Crew         *general.Crew `json:"crew"`          // Crew that took a part in production
	ReleasedYear int           `json:"released_year"` //Year the show was released in
	Revenue      float64       `json:"revenue"`       // Revenue which was received by the show
	Genre        string        `json:"genre"`         // A genre of the show
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
