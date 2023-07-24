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
	"github.com/google/uuid"
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
			apiTest.mockBehaviour(filmService, *db.NewFilm(generator.Generate(), apiTest.inputFilm.Title, apiTest.inputFilm.Crew, apiTest.inputFilm.ReleasedYear, apiTest.inputFilm.Revenue, apiTest.inputFilm.Genre))
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

func TestAggregatorServer_GetAllFilms(t *testing.T) {
	dbFilm := db.NewFilm(uuid.New(), "XXX", nil, 2004, 10002.99, "fantasy")
	type mockBehaviour func(s *mock_service.MockService)
	testTable := []struct {
		testName           string
		mockBehaviour      mockBehaviour
		expectedStatusCode int
		insertedFilm       db.Film
	}{
		{
			testName: "Not found",
			mockBehaviour: func(s *mock_service.MockService) {
				s.EXPECT().GetAll(context.Background()).Return([]api.Film{}, nil)
			},
			expectedStatusCode: 404,
		}, {
			testName: "Found",
			mockBehaviour: func(s *mock_service.MockService) {
				s.EXPECT().GetAll(context.Background()).Return([]api.Film{*api.NewFilm(*dbFilm)}, nil)
			},
			insertedFilm:       *dbFilm,
			expectedStatusCode: 200,
		}, {
			testName: "InternalError",
			mockBehaviour: func(s *mock_service.MockService) {
				s.EXPECT().GetAll(context.Background()).Return(nil, fmt.Errorf("error while scanning"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, apiTest := range testTable {
		t.Run(apiTest.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filmService := mock_service.NewMockService(ctrl)

			generator := &mocks.MockIDGenerator{}
			server := AggregatorServer{FilmsService: filmService, IDGenerator: generator}
			server.Router = http.NewServeMux()
			apiTest.mockBehaviour(filmService)
			server.Router.HandleFunc("/api/all/", server.GetAllFilms)
			w := httptest.NewRecorder()
			path := fmt.Sprintf("http://localhost:8080/api/all/")
			req := httptest.NewRequest("GET", path, nil)
			server.Router.ServeHTTP(w, req)

			assert.Equal(t, apiTest.expectedStatusCode, w.Code)
		})
	}
}

func TestAggregatorServer_DeleteFilm(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockService, request api.DeleteRequest)
	testTable := []struct {
		title              string
		deleteRequest      api.DeleteRequest
		expectedStatusCode int
		mockBehaviour      mockBehaviour
	}{
		{
			title:              "Deleted",
			deleteRequest:      api.DeleteRequest{Title: "XXX", Genre: "Fantasy", ReleasedYear: 2004},
			expectedStatusCode: 200,
			mockBehaviour: func(s *mock_service.MockService, request api.DeleteRequest) {
				s.EXPECT().Delete(context.Background(), request).Return(nil)
			},
		}, {
			title:              "Bad Request",
			deleteRequest:      api.DeleteRequest{},
			expectedStatusCode: 400,
			mockBehaviour: func(s *mock_service.MockService, request api.DeleteRequest) {
				s.EXPECT().Delete(context.Background(), request).Return(nil).AnyTimes()
			},
		}, {
			title:         "Error while deleting",
			deleteRequest: api.DeleteRequest{Title: "YYY", Genre: "fantasy", ReleasedYear: 1996},
			mockBehaviour: func(s *mock_service.MockService, request api.DeleteRequest) {
				s.EXPECT().Delete(context.Background(), request).Return(fmt.Errorf("error while deleting"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, apiTest := range testTable {
		t.Run(apiTest.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filmService := mock_service.NewMockService(ctrl)

			generator := &mocks.MockIDGenerator{}
			server := AggregatorServer{FilmsService: filmService, IDGenerator: generator}
			server.Router = http.NewServeMux()
			apiTest.mockBehaviour(filmService, apiTest.deleteRequest)
			server.Router.HandleFunc("/api/delete/", server.DeleteFilm)
			w := httptest.NewRecorder()
			urlPath := fmt.Sprintf("http://localhost:8080/api/delete/?title=%s&genre=%s&released_year=%d",
				apiTest.deleteRequest.Title,
				apiTest.deleteRequest.Genre,
				apiTest.deleteRequest.ReleasedYear)
			req := httptest.NewRequest("DELETE", urlPath, nil)
			server.Router.ServeHTTP(w, req)

			assert.Equal(t, apiTest.expectedStatusCode, w.Code)
		})
	}
}
