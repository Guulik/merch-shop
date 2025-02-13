package repository

import (
	"context"
	"merch/internal/lib/logger"
)

func (r *Repo) CheckUserByUsername(ctx context.Context, username string) (int, error) {
	//TODO: wrap sql with squirrel
	var (
		query = `
		SELECT id
		FROM users 
		WHERE username = $1;
`
		values = []any{username}
		userId int
	)
	err := r.dbPool.QueryRow(ctx, query, values...).Scan(&userId)

	if err != nil {
		return -1, logger.WrapError(ctx, err)
	}
	return userId, nil
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
	err := r.dbPool.QueryRow(ctx, query, values...).Scan(&userId)
	if err != nil {
		return 0, logger.WrapError(ctx, err)
	}

	return userId, nil
}
