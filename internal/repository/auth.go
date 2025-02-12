package repository

import (
	"context"
	"merch/internal/domain/model"
)

func (r *Repo) CheckUserByUsername(ctx context.Context, username string) (*model.UserAuth, error) {
	var (
		query = `
		SELECT id
		FROM users 
		WHERE username = $1;
`
		values = []any{username}
		user   *model.UserAuth
	)
	err := r.dbPool.Get(user, query, values...)
	if err != nil {
		//TODO: log error
		return nil, err
	}
	return user, nil
}

func (r *Repo) SaveUser(ctx context.Context, username string, password string) (int, error) {

	var (
		query = `
		INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING id;
	`
		values = []any{username, password}

		userId int
	)
	err := r.dbPool.QueryRow(query, values...).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
