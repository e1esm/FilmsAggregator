package main

import (
	"github.com/e1esm/FilmsAggregator/internal/server"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	currServer := configureServer()
	log.Fatal(http.ListenAndServe(cfg.Aggregator.Address+cfg.Aggregator.Port,
		currServer.Router))
}

func configureServer() *server.AggregatorServer {
	sb := server.NewBuilder()
	aggServ := sb.WithRouter(http.NewServeMux()).
		WithEndpoint("/api/add/", sb.Server.AddFilm).
		Build()
	return aggServ
}
