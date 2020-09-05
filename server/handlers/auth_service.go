package handlers

import (
	"VideoHub/database"
	"VideoHub/database/models"
	api "VideoHub/pb"
	"VideoHub/server/utils"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	jwtManager *utils.JWTManager
	userRepo database.AuthRepo
}

func (a *AuthService) SignIn(ctx context.Context, request *api.SignInRequest) (*api.SignInResponse, error) {
	user, err := a.userRepo.Get(ctx, request.Login)
	if err != nil {
		// TODO чек на ашибку notExists
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if user.Password != request.Password {
		return nil, status.Error(codes.PermissionDenied, "passwords not match")
	}

	jwt, err := a.jwtManager.Sign(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "jwt manger can not sign jwt")
	}

	return &api.SignInResponse{
		Token: jwt,
	}, nil
}

func (a *AuthService) SignOut(ctx context.Context, request *api.SignOutRequest) (*api.SignOutResponse, error) {
	if request.Password != "" && request.Password != request.PasswordConfirmation {
		return nil, status.Error(codes.InvalidArgument, "passwords must be same")
	}

	id, err := a.userRepo.Create(ctx, &models.User{
		Username: request.Username,
		Email: request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.SignOutResponse{
		Success: fmt.Sprintf("user created successfuly %d", id),
	}, nil
}

func NewAuthService(jwtManger *utils.JWTManager, userrepo database.AuthRepo) api.AuthServiceServer {
	return &AuthService{
		userRepo: userrepo,
		jwtManager: jwtManger,
	}
}

