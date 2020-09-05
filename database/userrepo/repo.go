package userrepo

import (
	"VideoHub/database"
	"VideoHub/database/models"
	"context"
	"errors"
)

type userRepository struct {
	sequence int64
	db       map[int64]*models.User
}

func (u *userRepository) Get(ctx context.Context, username string) (*models.User, error) {
	for _, v := range u.db {
		if v.Username == username || v.Email == username {
			return v, nil
		}
	}

	return nil, errors.New("user not found")
}

func (u *userRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	user.ID = u.sequence
	u.db[u.sequence] = user
	u.sequence++

	return user.ID, nil
}

func (u *userRepository) Update(ctx context.Context) (int64, error) {
	panic("implement me")
}

func (u *userRepository) Remove(ctx context.Context) {
	panic("implement me")
}

func New() database.AuthRepo {
	return &userRepository{
		sequence: 1,
		db: make(map[int64]*models.User),
	}
}
