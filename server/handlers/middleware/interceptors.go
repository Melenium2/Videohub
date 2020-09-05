package middleware

import (
	"VideoHub/server/utils"
	"encoding/json"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

type Middleware interface {
	HttpChain(http.Handler) http.Handler
	GrpcChain() []grpc.ServerOption
}

type Interceptors struct {
	jwtManger  *utils.JWTManager
	exceptions []string
	logger     *zap.Logger
}

func New(manger *utils.JWTManager, logger *zap.Logger, exceptions ...string) Middleware {
	return &Interceptors{
		jwtManger:  manger,
		exceptions: exceptions,
		logger: logger,
	}
}

//Выстроить chain из middleware for http и вернуть handler
func (in *Interceptors) HttpChain(next http.Handler) http.Handler {
	return in.httpAuth(next)
}

func (in *Interceptors) GrpcChain() []grpc.ServerOption {
	return append(AddLogging(in.logger, []grpc.ServerOption{}...), grpc.ChainUnaryInterceptor(
		UnaryInterceptor(in.grpcAuth),
		grpc_prometheus.UnaryServerInterceptor,
	))
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			writeError(w, http.StatusInternalServerError, err)
		}
	}
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeResponse(w, code, map[string]interface{}{
		"error": err.Error(),
	})
}

//https://medium.com/@amsokol.com/tutorial-part-3-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-739aac8f1d7e
//https://github.com/grpc-ecosystem/go-grpc-middleware
