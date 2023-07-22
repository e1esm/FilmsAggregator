package db

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"time"
)

type Film struct {
	ID           uuid.UUID    `json:"-" reindex:"-"`
	CacheID      int64        `json:"ID" reindex:"id,,pk"`
	Title        string       `json:"title" reindex:"title, tree"`
	Crew         general.Crew `json:"crew"`
	Revenue      float64      `json:"revenue" reindex:"revenue"`
	ReleasedYear int          `json:"released_year" reindex:"released_year"`
	CacheTime    time.Time    `json:"cache_time" reindex:"cache_time"`
	HashCode     [16]byte     `reindex:"hash"`
}

func NewFilm(film Film) *Film {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(film)
	if err != nil {
		return nil
	}
	return &Film{ID: film.ID,
		Title:        film.Title,
		Crew:         film.Crew,
		Revenue:      film.Revenue,
		CacheTime:    time.Now(),
		HashCode:     md5.Sum(b.Bytes()),
		ReleasedYear: film.ReleasedYear,
	}
}
