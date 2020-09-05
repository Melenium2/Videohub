package database

import (
	"VideoHub/database/models"
	"context"
)

type AuthRepo interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	Update(ctx context.Context) (int64, error)
	Remove(ctx context.Context)
	Get(ctx context.Context, username string) (*models.User, error)
}
