package main

import (
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/internal/server"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"log"
	"net/http"
)

func main() {

	cfg := config.NewConfig()
	currServer := configureServer()

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d",
		cfg.Aggregator.Address,
		cfg.Aggregator.Port),
		currServer.Router))
}

func configureServer() *server.AggregatorServer {
	sb := server.NewBuilder()
	aggServ := sb.WithRouter(http.NewServeMux()).
		WithEndpoint("/api/add/", sb.Server.AddFilm).
		Build()
	return aggServ
}

func configureRepositories() *repository.Repositories {

}

func configureService() *service.FilmsService {

}
