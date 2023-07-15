package server

import "net/http"

type AggregatorServer struct {
	router *http.ServeMux
}

type Builder struct {
	server *AggregatorServer
}

func NewBuilder() *Builder {
	builder := &Builder{server: &AggregatorServer{}}
	return builder
}

func (b *Builder) WithEndpoint(endpoint string, handlerFunc http.HandlerFunc) *Builder {
	b.server.router.HandleFunc(endpoint, handlerFunc)
	return b
}

func (b *Builder) WithRouter(mux *http.ServeMux) *Builder {
	b.server.router = mux
	return b
}

func (b *Builder) WithService() *Builder {
	//b.server.service = serviceName
	return b
}

func (b *Builder) Build() *AggregatorServer {
	return b.server
}
