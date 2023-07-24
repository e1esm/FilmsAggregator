package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/service"
	mock_service "github.com/e1esm/FilmsAggregator/internal/service/mocks"
	"github.com/e1esm/FilmsAggregator/utils/uuid/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAggregatorServer_AddFilm(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockService, film db.Film)
	testTable := []struct {
		name               string
		inputFilm          api.Film
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:               "OK",
			inputFilm:          api.Film{Title: "XXX", Revenue: 10000, ReleasedYear: 2004},
			expectedStatusCode: 200,
			mockBehaviour: func(s *mock_service.MockService, film db.Film) {
				s.EXPECT().Add(context.Background(), film).Return(*api.NewFilm(film), nil)
			},
		},
		{
			name:               "Fail",
			inputFilm:          api.Film{},
			expectedStatusCode: 400,
			mockBehaviour: func(s *mock_service.MockService, film db.Film) {
				s.EXPECT().Add(context.Background(), film).Return(api.Film{}, service.AlreadyExistsError)

			},
		},
	}

	for _, apiTest := range testTable {
		t.Run(apiTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filmService := mock_service.NewMockService(ctrl)

			generator := &mocks.MockIDGenerator{}
			apiTest.mockBehaviour(filmService, *db.NewFilm(generator.Generate(), apiTest.inputFilm.Title, apiTest.inputFilm.Crew, apiTest.inputFilm.ReleasedYear, apiTest.inputFilm.Revenue))
			server := AggregatorServer{FilmsService: filmService, IDGenerator: generator}
			server.Router = http.NewServeMux()

			server.Router.HandleFunc("/api/add/", server.AddFilm)
			content, _ := json.Marshal(apiTest.inputFilm)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "http://localhost:8080/api/add/", bytes.NewBufferString(string(content)))
			server.Router.ServeHTTP(w, req)

			assert.Equal(t, apiTest.expectedStatusCode, w.Code)
		})
	}
}

func TestAggregatorServer_GetFilms(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockService, name string)
	testTable := []struct {
		testName           string
		filmName           string
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			testName: "OK",
			filmName: "XXX",
			mockBehaviour: func(s *mock_service.MockService, name string) {
				s.EXPECT().Get(context.Background(), name).Return([]*api.Film{}, nil)
			},
			expectedStatusCode: 200,
		},
	}
	for _, apiTest := range testTable {
		t.Run(apiTest.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filmService := mock_service.NewMockService(ctrl)

			generator := &mocks.MockIDGenerator{}
			apiTest.mockBehaviour(filmService, apiTest.filmName)
			server := AggregatorServer{FilmsService: filmService, IDGenerator: generator}
			server.Router = http.NewServeMux()

			server.Router.HandleFunc("/api/get", server.GetFilms)
			w := httptest.NewRecorder()
			path := fmt.Sprintf("http://localhost:8080/api/get?name=%s", apiTest.filmName)
			req := httptest.NewRequest("GET", path, nil)
			server.Router.ServeHTTP(w, req)

			assert.Equal(t, apiTest.expectedStatusCode, w.Code)
		})
	}
}
