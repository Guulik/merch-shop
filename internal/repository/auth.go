package repository

import (
	"context"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
)

func (r *Repo) CheckUserByUsername(ctx context.Context, username string) (*model.UserAuth, error) {
	var (
		query = `
		SELECT id, username, password_hash
		FROM users 
		WHERE username = $1;
`
		values = []any{username}
		user   = &model.UserAuth{Id: -1}
	)
	err := r.dbPool.QueryRow(ctx, query, values...).Scan(&user.Id, &user.Username, &user.PasswordDb)

	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return user, nil
}

func (r *Repo) SaveUser(ctx context.Context, username string, password string) (int, error) {

	var (
		query = `
		INSERT INTO users (username, password_hash, coins) 
		VALUES ($1, $2, $3) 
		RETURNING id;
	`
		values = []any{username, password, 1000}

		userId int
	)
	err := r.dbPool.QueryRow(ctx, query, values...).Scan(&userId)
	if err != nil {
		return 0, logger.WrapError(ctx, err)
	}

	return userId, nil
}
