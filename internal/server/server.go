package server

import (
	"github.com/e1esm/FilmsAggregator/internal/service"
	"net/http"
)

type AggregatorServer struct {
	Router       *http.ServeMux
	FilmsService service.Service
}

type Builder struct {
	Server *AggregatorServer
}

func NewBuilder() *Builder {
	builder := &Builder{Server: &AggregatorServer{}}
	return builder
}

func (b *Builder) WithEndpoint(endpoint string, handlerFunc http.HandlerFunc) *Builder {
	b.Server.Router.HandleFunc(endpoint, handlerFunc)
	return b
}

func (b *Builder) WithRouter(mux *http.ServeMux) *Builder {
	b.Server.Router = mux
	return b
}

func (b *Builder) WithService(service service.Service) *Builder {
	b.Server.FilmsService = service
	return b
}

func (b *Builder) Build() *AggregatorServer {
	return b.Server
}
