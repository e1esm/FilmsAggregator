package repository

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models"
)

type Repository interface {
	Add(context.Context, *models.Film) (models.Film, error)
}

type Cache interface {
	Repository
}

type Repositories struct {
	CacheRepo Cache
	MainRepo  Repository
}

type RepositoriesBuilder struct {
	Repositories Repositories
}

func NewRepositoriesBuilder() *RepositoriesBuilder {
	reposBuilder := RepositoriesBuilder{Repositories: Repositories{}}
	return &reposBuilder
}

func (rb *RepositoriesBuilder) WithMainRepo(repository Repository) *RepositoriesBuilder {
	rb.Repositories.MainRepo = repository
	return rb
}

func (rb *RepositoriesBuilder) WithCacheRepo(repository Cache) *RepositoriesBuilder {
	rb.Repositories.CacheRepo = repository
	return rb
}

func (rb *RepositoriesBuilder) Build() *Repositories {
	return &rb.Repositories
}
