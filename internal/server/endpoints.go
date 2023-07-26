package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

// AddFilm godoc
// @Summary Add film to the DB
// @Description Based on the body of POST request add film to the DB
// @Tags film
// @Accept json
// @Produce json
// @Param film body api.Film true "film model"
// @Success 200 {object} api.Film
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /api/add/ [post]
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

	w.Header().Set("Content-Type", "application/json")

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
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// GetFilms godoc
// @Summary Get films with the specified name (there can be more than 1 film with the same name)
// @Description Get all films with the specified name.
// @Param name query string true "film title"
// @Tags films
// @Produce json
// @Success 200 {array} api.Film
// @Failure 400
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /api/get/ [get]
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

// DeleteFilm godoc
// @Summary Delete film
// @Description Delete film from both cache and main repositories based on the user's provided filters
// @Tags film
// @Param title query string true "Film title"
// @Param genre query string true "Film genre"
// @Param released_year query string true "Film release date"
// @Produce json
// @Success 200 {object} api.DeleteRequest
// @Failure 400
// @Failure 405
// @Failure 500
// @Router /api/delete/ [delete]
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

	err = s.FilmsService.Delete(context.Background(), deleteRequest)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%v", deleteRequest))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content, err := json.Marshal(deleteRequest)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// GetAllFilms godoc
// @Summary Get all the films from the DB
// @Description Get every available film from the DB
// @Tags films
// @Produce json
// @Success 200 {array} api.Film
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /api/all/ [get]
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
	if len(films) == 0 {
		w.WriteHeader(http.StatusNotFound)
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

// FindFilmsByActor godoc
// @Summary Films actor took a part in
// @Description Get all the films which a certain actor was shot in
// @Param actor query string true "Actor filter"
// @Tags actor
// @Produce json
// @Success 200 {array} api.Film
// @Failure 400
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /api/actor/films/ [get]
func (s *AggregatorServer) FindFilmsByActor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	actorName := r.URL.Query().Get("actor")
	if actorName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	films, err := s.FilmsService.GetByActor(r.Context(), actorName)
	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	content, err := json.Marshal(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// FindFilmsByProducer godoc
// @Summary Films, which were produced by the specified person
// @Description Get all the films that'd been produced by a specified producer.
// @Param producer query string true "Producer filter"
// @Tags producer
// @Produce json
// @Success 200 {array} api.Film
// @Failure 400
// @Failure 404
// @Failure 405
// @Failure 500
// @Router /api/producer/films/ [get]
func (s *AggregatorServer) FindFilmsByProducer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	producerName := r.URL.Query().Get("producer")
	if producerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	films, err := s.FilmsService.GetByProducer(r.Context(), producerName)
	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	content, err := json.Marshal(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
