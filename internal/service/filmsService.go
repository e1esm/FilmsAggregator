package service

import "github.com/e1esm/FilmsAggregator/internal/repository"

type Service interface {
}

type FilmsService struct {
	Repositories *repository.Repositories
}

func NewService(repositories *repository.Repositories) *FilmsService {
	return &FilmsService{Repositories: repositories}
}
