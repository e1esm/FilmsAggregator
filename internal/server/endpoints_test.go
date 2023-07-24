package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	mock_service "github.com/e1esm/FilmsAggregator/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockIDGenerator struct {
}

func (m *MockIDGenerator) Generate() uuid.UUID {
	return uuid.UUID{}
}
func (m *MockIDGenerator) GenerateUUIDs(film db.Film) *db.Film {
	return &film
}

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
				s.EXPECT().Add(context.Background(), film).Return(*api.NewFilm(film), nil).AnyTimes()
			},
		},
	}

	for _, apiTest := range testTable {
		t.Run(apiTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filmService := mock_service.NewMockService(ctrl)

			generator := &MockIDGenerator{}
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
