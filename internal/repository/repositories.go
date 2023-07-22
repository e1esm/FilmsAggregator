package repository

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
)

type Repository interface {
	Add(context.Context, *api.Film) (api.Film, error)
	FindByName(ctx context.Context, name string) ([]*api.Film, error)
}

type Cache interface {
	Repository
	Delete(context.Context, string) (api.Film, error)
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
