package server

import (
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (s *AggregatorServer) AddFilm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Logger.Error("Another method was used on this URL",
			zap.String("url", r.URL.String()))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Logger.Error("Couldn't have read the body of the request",
			zap.String("err", err.Error()),
			zap.String("url", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	receivedFilm := &models.Film{}
	err = json.Unmarshal(content, receivedFilm)
	if err != nil {
		logger.Logger.Error("Couldn't have unmarshalled request's body",
			zap.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
