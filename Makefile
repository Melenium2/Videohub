
GOBIN=./cmd/main

## build: Build go binary
build:
	go build -o $(GOBIN)

## run: Run go server
.PHONY: run
run:
	$(GOBIN)

## get: Run go get missing dependencies
get:
	go get ./...

## deploy: Run commands to deploy app to container
deploy:
	make get
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
          -ldflags='-w -s -extldflags "-static"' -a \
          -o /go/bin/main .

## generate go protobuf files. you need to specify your own path to utilities
gen:
	protoc -I./proto \
                -I ${GOPATH}/bin \
                -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.6/third_party/googleapis \
                -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.6 \
                --grpc-gateway_out=logtostderr=true:./pb \
                --swagger_out=allow_merge=true,merge_file_name=api:. \
                --go_out=plugins=grpc:./pb ./proto/authserver.proto ./proto/swagger.proto ./proto/simple.proto

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := help