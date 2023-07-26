package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/logger"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	valid_cache_time_minutes = 15
)

var (
	AlreadyExistsError = errors.New("film already exists")
)

//go:generate mockgen -source=filmsService.go -destination=mocks/mock.go

type Service interface {
	Add(context.Context, db.Film) (api.Film, error)
	Get(ctx context.Context, name string) ([]*api.Film, error)
	Delete(ctx context.Context, requestedFilm api.DeleteRequest) error
	GetAll(ctx context.Context) ([]api.Film, error)
	GetByActor(ctx context.Context, name string) ([]api.Film, error)
	GetByProducer(ctx context.Context, name string) ([]api.Film, error)
}

type FilmsService struct {
	Repositories   *repository.Repositories
	ExpirationTime int
}

func NewService(repositories *repository.Repositories, config *config.Config, generator uuid.Generator) *FilmsService {
	expirationTime, err := strconv.Atoi(config.Reindexer.CacheTime)
	if err != nil {
		logger.Logger.Error(err.Error())
		expirationTime = valid_cache_time_minutes
	}
	return &FilmsService{Repositories: repositories, ExpirationTime: expirationTime}
}

func (fs *FilmsService) Add(ctx context.Context, film db.Film) (api.Film, error) {
	film.CacheTime = time.Now()
	doesExist := fs.Repositories.MainRepo.Verify(ctx, &film)
	if doesExist {
		return api.Film{}, AlreadyExistsError
	}
	_, err := fs.Repositories.CacheRepo.Add(ctx, film)
	if err != nil {
		logger.Logger.Error(err.Error(), zap.String("film", film.Title))
	}
	inserted, err := fs.Repositories.MainRepo.Add(ctx, film)
	if err != nil {
		return api.Film{}, err
	}
	return *api.NewFilm(inserted), nil
}

func (fs *FilmsService) Get(ctx context.Context, name string) ([]*api.Film, error) {
	received, err := fs.Repositories.CacheRepo.FindByName(ctx, name)
	apiFilms := make([]*api.Film, 0)
	if err != nil {
		logger.Logger.Error("Couldn't have retrieved films from cache", zap.String("err", err.Error()))
	}
	current := time.Now()
	isUpToDate := true
	for i := 0; i < len(received); i++ {
		if int(current.Sub(received[i].CacheTime).Minutes()) > fs.ExpirationTime {
			isUpToDate = false
		}
	}

	if isUpToDate && len(received) > 0 {

		for i := 0; i < len(received); i++ {
			apiFilms = append(apiFilms, api.NewFilm(*received[i]))
		}
		logger.Logger.Info(fmt.Sprintf("%v", received))
		return apiFilms, nil
	} else {
		for i := 0; i < len(received); i++ {
			err = fs.Repositories.CacheRepo.Delete(ctx, api.DeleteRequest{received[i].Genre, received[i].Title, received[i].ReleasedYear})
			if err != nil {
				logger.Logger.Error("Couldn't have deleted film from cache",
					zap.String("err", err.Error()),
					zap.String("film", received[i].Title))
			}
		}
	}

	received, err = fs.Repositories.MainRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(received); i++ {
		apiFilms = append(apiFilms, api.NewFilm(*received[i]))
	}
	return apiFilms, nil
}

func (fs *FilmsService) Delete(ctx context.Context, request api.DeleteRequest) error {

	cacheCtx := context.WithValue(ctx, "request", request)
	CacheErr := fs.Repositories.CacheRepo.DeleteCachedWithCtx(cacheCtx)
	if CacheErr != nil {
		logger.Logger.Error(CacheErr.Error())
	}

	MainErr := fs.Repositories.MainRepo.Delete(ctx, request)
	if MainErr != nil {
		logger.Logger.Error(MainErr.Error())
		return MainErr
	}
	return nil
}

func (fs *FilmsService) GetAll(ctx context.Context) ([]api.Film, error) {
	films, err := fs.Repositories.MainRepo.FindAll(ctx)
	filmsAPI := make([]api.Film, len(films))
	for i := 0; i < len(films); i++ {
		filmsAPI[i] = *api.NewFilm(films[i])
	}
	if err != nil {
		return filmsAPI, err
	}
	return filmsAPI, nil
}

func (fs *FilmsService) GetByActor(ctx context.Context, name string) ([]api.Film, error) {
	films, err := fs.Repositories.MainRepo.FindFilmsByActor(ctx, name)
	if err != nil {
		return nil, err
	}
	apiFilms := make([]api.Film, len(films))
	for i := 0; i < len(films); i++ {
		apiFilms[i] = *api.NewFilm(films[i])
	}
	return apiFilms, nil
}

func (fs *FilmsService) GetByProducer(ctx context.Context, name string) ([]api.Film, error) {
	films, err := fs.Repositories.MainRepo.FindFilmsByProducer(ctx, name)
	if err != nil {
		return nil, err
	}
	apiFilms := make([]api.Film, len(films))
	for i := 0; i < len(films); i++ {
		apiFilms[i] = *api.NewFilm(films[i])
	}
	return apiFilms, nil
}
