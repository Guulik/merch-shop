package repository

import (
	"context"
	"merch/internal/domain/model"
)

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (model.UserAuth, error) {
	//TODO implement me
	var (
		user model.UserAuth
	)
	return user, nil
}

func (r *Repo) SaveUser(ctx context.Context, username string, password string) error {
	//TODO implement me
	return nil
}
