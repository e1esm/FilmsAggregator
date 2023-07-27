package server

import (
	"github.com/e1esm/FilmsAggregator/internal/service"
	"github.com/e1esm/FilmsAggregator/utils/uuid"
	"net/http"
)

type AggregatorServer struct {
	Router       *http.ServeMux
	FilmsService service.Service
	IDGenerator  uuid.Generator
	AuthService  service.AuthorizationService
}

type Builder struct {
	Server *AggregatorServer
}

func NewBuilder() *Builder {
	builder := &Builder{Server: &AggregatorServer{}}
	return builder
}

func (b *Builder) WithIDGenerator(idgen uuid.Generator) *Builder {
	b.Server.IDGenerator = idgen
	return b
}

func (b *Builder) WithEndpoint(endpoint string, handlerFunc http.HandlerFunc) *Builder {
	b.Server.Router.HandleFunc(endpoint, handlerFunc)
	return b
}

func (b *Builder) WithProtectedEndpoint(endpoint string, handler http.Handler) *Builder {
	b.Server.Router.Handle(endpoint, handler)
	return b
}

func (b *Builder) WithRouter(mux *http.ServeMux) *Builder {
	b.Server.Router = mux
	return b
}

func (b *Builder) WithFilmsService(service service.Service) *Builder {
	b.Server.FilmsService = service
	return b
}

func (b *Builder) WithAuthenticationService(service service.AuthorizationService) *Builder {
	b.Server.AuthService = service
	return b
}

func (b *Builder) Build() *AggregatorServer {
	return b.Server
}
