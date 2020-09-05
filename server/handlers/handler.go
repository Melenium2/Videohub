package handlers

import (
	"VideoHub/database/userrepo"
	api "VideoHub/pb"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Handler interface {
	Setup(ctx context.Context, server *grpc.Server) []HttpEndpointUnit
}

type HttpEndpointFunc func(uint, *runtime.ServeMux, ...grpc.DialOption) error

type HttpEndpointUnit struct {
	Path     string
	Secured  bool
	Endpoint HttpEndpointFunc
}

type serviceHandler struct {
	*Config
}

func (s *serviceHandler) Setup(ctx context.Context, server *grpc.Server) []HttpEndpointUnit {
	httpEndpoints := make([]HttpEndpointUnit, 0)
	api.RegisterAuthServiceServer(server, NewAuthService(s.JwtManager, userrepo.New()))
	api.RegisterSimpleServiceServer(server, NewSimpleService())

	authUnit := HttpEndpointUnit{
		Path:    "/auth",
		Secured: false,
		Endpoint: func(port uint, mux *runtime.ServeMux, options ...grpc.DialOption) error {
			return api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", port), options)
		},
	}
	simpleUnit := HttpEndpointUnit{
		Path:    "/simple",
		Secured: true,
		Endpoint: func(port uint, mux *runtime.ServeMux, options ...grpc.DialOption) error {
			return api.RegisterSimpleServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", port), options)
		},
	}

	/*
		TODO
		Сделать middleware чтоб не пускало без токена
	*/

	return append(
		httpEndpoints,
		authUnit,
		simpleUnit,
	)
}

func New(config *Config) Handler {
	return &serviceHandler{config}
}
