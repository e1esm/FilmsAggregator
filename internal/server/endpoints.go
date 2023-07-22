package server

import (
	"context"
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func (s *AggregatorServer) AddFilm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	receivedFilm := &api.Film{}
	err = json.Unmarshal(content, receivedFilm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dtoFilm := db.NewFilm(receivedFilm.ID, receivedFilm.Title, receivedFilm.Crew, receivedFilm.ReleasedYear, receivedFilm.Revenue)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	insertedFilm, err := s.FilmsService.Add(ctx, dtoFilm)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(insertedFilm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *AggregatorServer) GetFilms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query()
	name := params.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	films, err := s.FilmsService.Get(ctx, name)
	if err != nil {
		logger.Logger.Error(err.Error(), zap.String("req", name))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bytes, err := json.Marshal(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
