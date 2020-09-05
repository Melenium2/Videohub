package handlers

import (
	api "VideoHub/pb"
	"context"
)

type SimpleService struct {

}

func (s SimpleService) DoIt(ctx context.Context, request *api.SimpleRequest) (*api.SimpleResponse, error) {
	return &api.SimpleResponse{
		Obj: request.Obj,
	}, nil
}

func NewSimpleService() *SimpleService {
	return &SimpleService{}
}