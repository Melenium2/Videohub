package main

import (
	"VideoHub/server"
	"VideoHub/server/handlers"
	"VideoHub/server/handlers/middleware"
	"VideoHub/server/logger"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	if err := logger.Init(0); err != nil {
		log.Fatal(err)
	}

	handlersConfig := handlers.NewConfig()
	endpoints := handlers.New(handlersConfig)

	interceptors := middleware.New(handlersConfig.JwtManager, logger.Log, "signin", "signout")
	config := server.NewConfig(9876, 9880, interceptors)

	s, err := server.New(ctx, config, endpoints)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
