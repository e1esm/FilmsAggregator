package repository

import (
	"github.com/e1esm/FilmsAggregator/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Add(film models.Film) uuid.UUID
}

type Repositories struct {
	CacheRepo *Repository
	MainRepo  *Repository
}

type RepositoriesBuilder struct {
	Repositories Repositories
}

func NewRepositoriesBuilder() *RepositoriesBuilder {
	reposBuilder := RepositoriesBuilder{Repositories: Repositories{}}
	return &reposBuilder
}

func (rb *RepositoriesBuilder) WithMainRepo(repository *Repository) *RepositoriesBuilder {
	rb.Repositories.MainRepo = repository
	return rb
}

func (rb *RepositoriesBuilder) WithCacheRepo(repository *Repository) *RepositoriesBuilder {
	rb.Repositories.CacheRepo = repository
	return rb
}

func (rb *RepositoriesBuilder) Build() *Repositories {
	return &rb.Repositories
}
