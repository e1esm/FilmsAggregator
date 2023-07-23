package repository

import (
	"context"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
)

type Repository interface {
	Add(context.Context, *db.Film) (db.Film, error)
	FindByName(ctx context.Context, name string) ([]*db.Film, error)
	Verify(ctx context.Context, film *db.Film) bool
}

type Cache interface {
	Repository
	Delete(context.Context, string) (db.Film, error)
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
