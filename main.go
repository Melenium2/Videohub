package main

import (
	"VideoHub/server"
	"VideoHub/server/handlers"
	"VideoHub/server/handlers/middleware"
	"VideoHub/server/logger"
	"context"
	"log"
)



/**
	Почитать про то как писать swagger файлы
	Переписать Видеохаб для того чтобы нормально раздлить логику grpc и http
		сейчас оно все смешано
	TODO
	MAIN:
	Поменять makefile таким образом чтобы каждый протофай генерил свой swagger.json
	Каждый протофайл разделить по разным пекеджем

	HLS:
	Стримить видео и разбивать его на сегменты на лету
		https://github.com/TechMaster/gohls/blob/master/hls/encoder.go
		https://github.com/hanhailong/golang-hls/blob/master/http_util.go
		https://dougsillars.com/2017/10/26/how-hls-adaptive-bitrate-works/
		https://www.rohitmundra.com/video-streaming-server
	Сделать алоад видео
		* Сервис
		* Контроллер
		* Тесты
	Сделать шаринг видео в бразузер
	Сделать кастомный битрейт(юзер может выбирать качество видео)
 */
func main() {
	ctx := context.Background()
	if err := logger.Init(0); err != nil {
		log.Fatal(err)
	}

	handlersConfig := handlers.NewConfig()
	endpoints := handlers.New(handlersConfig)

	// Замутить тут фабричный метод с конфигурацией на лету
	interceptors := middleware.New(handlersConfig.JwtManager, logger.Log, "signin", "signout")
	// Замутить тут фабричный метод с конфигурацией на лету
	config := server.NewConfig(9876, 9880, interceptors)

	s, err := server.New(ctx, config, endpoints)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
