package main

import (
	"fmt"
	_ "github.com/e1esm/FilmsAggregator/docs"
	"github.com/e1esm/FilmsAggregator/internal/repository"
	"github.com/e1esm/FilmsAggregator/internal/repository/authentication"
	"github.com/e1esm/FilmsAggregator/internal/repository/postgres"
	"github.com/e1esm/FilmsAggregator/internal/repository/reindexer"
	"github.com/e1esm/FilmsAggregator/internal/server"
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/config"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

// @title Films Aggregator
// @version 1.0
// @description API Server for Films Aggregator application
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.NewConfig()
	currServer := configureServer(configureService(cfg, &uuid.UUIDGenerator{}), service.NewAuthService(authentication.NewAuthRepository(*cfg)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d",
		cfg.Aggregator.Address,
		cfg.Aggregator.Port),
		currServer.Router))
}

func configureServer(service service.Service, authService service.AuthorizationService) *server.AggregatorServer {
	sb := server.NewBuilder()
	aggServ := sb.WithRouter(http.NewServeMux()).
		WithProtectedEndpoint("/api/add/", sb.Server.UserIdentity(http.HandlerFunc(sb.Server.AddFilm))).
		WithEndpoint("/api/get/", sb.Server.GetFilms).
		WithProtectedEndpoint("/api/delete/", sb.Server.UserIdentity(http.HandlerFunc(sb.Server.DeleteFilm))).
		WithEndpoint("/api/all/", sb.Server.GetAllFilms).
		WithEndpoint("/api/actor/films/", sb.Server.FindFilmsByActor).
		WithEndpoint("/api/producer/films/", sb.Server.FindFilmsByProducer).
		WithEndpoint("/api/signin/", sb.Server.SignIn).
		WithEndpoint("/api/signup/", sb.Server.SignUp).
		WithEndpoint("/swagger/", httpSwagger.WrapHandler).
		WithFilmsService(service).
		WithAuthenticationService(authService).
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
