package api

import (
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
