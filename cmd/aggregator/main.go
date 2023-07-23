package main

import (
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/internal/repository/postgres"
	"github.com/e1esm/FilmsAggregator/internal/repository/reindexer"
	"github.com/e1esm/FilmsAggregator/internal/server"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	currServer := configureServer(configureService(cfg))
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
		WithService(service).
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

func configureService(config *config.Config) *service.FilmsService {
	return service.NewService(configureRepositories(config), config)
}
