package db

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"time"
)

type Film struct {
	ID           uuid.UUID    `json:"-" reindex:"-"`
	CacheID      int64        `json:"ID" reindex:"id,,pk"`
	Title        string       `json:"title" reindex:"title,tree"`
	Crew         general.Crew `json:"crew"`
	Revenue      float64      `json:"revenue" reindex:"revenue"`
	ReleasedYear int          `json:"released_year" reindex:"released_year"`
	CacheTime    time.Time    `json:"cache_time" reindex:"cache_time"`
	HashCode     string       `reindex:"hash"`
}

func NewFilm(id uuid.UUID, title string, crew general.Crew, releasedYear int, revenue float64) *Film {
	film := &Film{ID: id, Title: title, Crew: crew, ReleasedYear: releasedYear, Revenue: revenue}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(film)
	if err != nil {
		return nil
	}
	film.HashCode = fmt.Sprintf("%v", md5.Sum(b.Bytes()))
	return film
}
