package server

import "net/http"

type AggregatorServer struct {
	Router *http.ServeMux
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

func (b *Builder) WithService() *Builder {
	//b.server.service = serviceName
	return b
}

func (b *Builder) Build() *AggregatorServer {
	return b.Server
}
