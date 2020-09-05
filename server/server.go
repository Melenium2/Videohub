package server

import (
	"VideoHub/server/handlers"
	"VideoHub/server/handlers/middleware"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Run() error
	Terminate()
}

type Config struct {
	port       uint
	httpPort   uint
	middleware middleware.Middleware
}

func NewConfig(port, httpPort uint, middleware middleware.Middleware) *Config {
	return &Config{
		port:       port,
		httpPort:   httpPort,
		middleware: middleware,
	}
}

type GrpcServer struct {
	ctx    context.Context
	config *Config
	errCh  chan error

	listener  net.Listener
	server    *grpc.Server
	endpoints handlers.Handler
	mux       *http.ServeMux
}

func (g *GrpcServer) Run() error {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		g.errCh <- g.server.Serve(g.listener)
	}()
	log.Printf("grpc server started on port %d", g.config.port)

	go func() {
		g.errCh <- http.ListenAndServe(fmt.Sprintf(":%d", g.config.httpPort), g.mux)
	}()
	log.Printf("server http proxy is now listening port %d", g.config.httpPort)

	select {
	case <-interrupt:
		log.Print("Server interrupted")
		g.Terminate()
		return nil
	case e := <-g.errCh:
		log.Printf("Error: %s", e)
		g.Terminate()
		return e
	}
}

func (g *GrpcServer) Terminate() {
	defer func() {
		r := recover()
		if r != nil {
			msg := "Server recover from panic"
			switch r.(type) {
			case string:
				log.Printf("Error: %s, trace: %s", msg, r)
			case error:
				log.Printf("Error: %s, trace: %s", msg, r)
			default:
				log.Printf("Error: %s", msg)
			}
		}

		return
	}()

	g.server.Stop()

	return
}

func New(ctx context.Context, config *Config, handlers handlers.Handler) (Server, error) {
	var err error

	if config == nil {
		return nil, errorConfigIsNil()
	}

	if config.port == 0 {
		return nil, errorPortIsEmpty()
	}

	server := &GrpcServer{
		config:    config,
		errCh:     make(chan error, 1),
		endpoints: handlers,
	}

	server.ctx = ctx
	if server.ctx == nil {
		server.ctx = context.Background()
	}

	{
		server.server = grpc.NewServer(server.config.middleware.GrpcChain()...)
	}
	log.Print("grpc server created")
	{
		server.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.port))
		if err != nil {
			return nil, err
		}
	}
	log.Print("listener created")
	if server.endpoints != nil {
		endpointsFuncs := server.endpoints.Setup(server.ctx, server.server)
		log.Print("create runtime mux")
		muxOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true})

		proxy := runtime.NewServeMux(muxOptions)

		options := []grpc.DialOption{grpc.WithInsecure()}
		for _, s := range endpointsFuncs {
			if err := s.Endpoint(server.config.port, proxy, options...); err != nil {
				return nil, err
			}
			log.Printf("%s endpoint done", s.Path)
		}
		mux := http.NewServeMux()
		mux.Handle("/", proxy)
		mux.Handle("/metrics", promhttp.Handler())
		server.mux = mux
	}
	log.Print("handler initialized")
	log.Print("server configured")
	return server, nil
}
