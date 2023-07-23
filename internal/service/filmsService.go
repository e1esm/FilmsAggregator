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

type Service interface {
	Add(context.Context, *db.Film) (api.Film, error)
	Get(ctx context.Context, name string) ([]*api.Film, error)
}

type FilmsService struct {
	Repositories   *repository.Repositories
	ExpirationTime int
}

func NewService(repositories *repository.Repositories, config *config.Config) *FilmsService {
	expirationTime, err := strconv.Atoi(config.Reindexer.CacheTime)
	if err != nil {
		logger.Logger.Error(err.Error())
		expirationTime = valid_cache_time_minutes
	}
	return &FilmsService{Repositories: repositories, ExpirationTime: expirationTime}
}

func (fs *FilmsService) Add(ctx context.Context, film *db.Film) (api.Film, error) {
	film.CacheTime = time.Now()
	doesExist := fs.Repositories.MainRepo.Verify(ctx, film)
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
			_, err = fs.Repositories.CacheRepo.Delete(ctx, received[i].Title)
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
