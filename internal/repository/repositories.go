package repository

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
)

type ScrapperRepository interface {
	FindAll(context.Context) ([]db.Film, error)
}

type CompleteRepository interface {
	ScrapperRepository
	MainRepository
}

type MainRepository interface {
	Add(context.Context, db.Film) (db.Film, error)
	FindByName(ctx context.Context, name string) ([]*db.Film, error)
	Verify(ctx context.Context, film *db.Film) bool
	Delete(context.Context, api.DeleteRequest) error
}

type Cache interface {
	MainRepository
}

type Repositories struct {
	CacheRepo Cache
	MainRepo  CompleteRepository
}

type RepositoriesBuilder struct {
	Repositories Repositories
}

func NewRepositoriesBuilder() *RepositoriesBuilder {
	reposBuilder := RepositoriesBuilder{Repositories: Repositories{}}
	return &reposBuilder
}

func (rb *RepositoriesBuilder) WithMainRepo(repository CompleteRepository) *RepositoriesBuilder {
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
