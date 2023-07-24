package main

import (
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/internal/repository/postgres"
	"github.com/e1esm/FilmsAggregator/internal/repository/reindexer"
	"github.com/e1esm/FilmsAggregator/internal/server"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	currServer := configureServer(configureService(cfg, &uuid.UUIDGenerator{}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d",
		cfg.Aggregator.Address,
		cfg.Aggregator.Port),
		currServer.Router))
}

func configureServer(service service.Service) *server.AggregatorServer {
	sb := server.NewBuilder()
	aggServ := sb.WithRouter(http.NewServeMux()).
		WithEndpoint("/api/add/", sb.Server.AddFilm).
		WithEndpoint("/api/get/", sb.Server.GetFilms).
		WithEndpoint("/api/delete/", sb.Server.DeleteFilm).
		WithEndpoint("/api/all/", sb.Server.GetAllFilms).
		WithService(service).
		WithIDGenerator(&uuid.UUIDGenerator{}).
		Build()
	return aggServ
}

func configureRepositories(config *config.Config) *repository.Repositories {
	mainRepo := repository.NewRepositoriesBuilder().
		WithMainRepo(postgres.NewFilmsRepository(*config, postgres.NewTransactionManager())).
		WithCacheRepo(reindexer.NewCacheRepository(*config)).
		Build()

	return mainRepo
}

func configureService(config *config.Config, generator uuid.Generator) *service.FilmsService {
	return service.NewService(configureRepositories(config), config, generator)
}
