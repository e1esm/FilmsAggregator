package server

import (
	"context"
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
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
	dtoFilm := db.NewFilm(s.IDGenerator.Generate(), receivedFilm.Title, receivedFilm.Crew, receivedFilm.ReleasedYear, receivedFilm.Revenue, receivedFilm.Genre)
	dtoFilm = s.IDGenerator.GenerateUUIDs(*dtoFilm)
	insertedFilm, err := s.FilmsService.Add(context.Background(), *dtoFilm)
	switch {
	case err == service.AlreadyExistsError:
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	case err != nil:
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

	films, err := s.FilmsService.Get(context.Background(), name)
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

func (s *AggregatorServer) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var (
		title             = r.URL.Query().Get("title")
		genre             = r.URL.Query().Get("genre")
		releasedYear, err = strconv.Atoi(r.URL.Query().Get("released_year"))
	)
	if title == "" || genre == "" || releasedYear == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteRequest := api.DeleteRequest{Title: title, Genre: genre, ReleasedYear: releasedYear}

	err = s.FilmsService.Delete(r.Context(), deleteRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content, err := json.Marshal(deleteRequest)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func (s *AggregatorServer) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	films, err := s.FilmsService.GetAll(r.Context())
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := json.Marshal(films)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
