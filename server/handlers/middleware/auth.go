package middleware

import (
	"VideoHub/server/utils"
	"context"
	"errors"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

const userKey uint = iota

var (
	errNotAuthorized = errors.New("not authorized")
)

type AuthFunc func(context.Context, string) (context.Context, error)

func UnaryInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		newCtx, err = authFunc(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

func (in *Interceptors) authorization(token string) (*utils.UserClaims, error) {
	user, err := in.jwtManger.Verify(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (in *Interceptors) grpcAuth(ctx context.Context, fullPath string) (context.Context, error) {
	for _, ex := range in.exceptions {
		fullPath = strings.ToLower(fullPath)
		if strings.Contains(fullPath, ex) {
			return ctx, nil
		}
	}

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	claims, err := in.authorization(token)
	if err != nil {
		return nil, err
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", claims)
	newCtx := context.WithValue(ctx, userKey, claims)

	return newCtx, nil
}

func (in *Interceptors) httpAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, ex := range in.exceptions {
			path := strings.ToLower(r.RequestURI)
			if strings.Contains(path, ex) {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := strings.Split(r.Header.Get("authorization"), "Bearer ")
		if len(authHeader) != 2 {
			writeError(w, http.StatusUnauthorized, errNotAuthorized)
			return
		}

		claims, err := in.authorization(authHeader[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, errNotAuthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), userKey, claims)
		r = r.WithContext(newCtx)

		next.ServeHTTP(w, r)
	})
}