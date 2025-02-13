package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"merch/internal/lib/logger"
)

func (r *Repo) TransferCoins(ctx context.Context, fromUserId int, toUserId int, coinAmount int) error {
	//TODO: wrap sql with squirrel
	var (
		subtractQuery = `
	UPDATE users 
	SET coins = coins - $1
	WHERE id = $2 AND coins >= $1
`
		addQuery = `
	UPDATE users 
	SET coins = coins + $1
	WHERE id = $2
`

		insertTransactionQuery = `
	INSERT INTO transactions (from_user, to_user, amount) 
	VALUES ($1, $2, $3)
`

		subtractQueryValues    = []any{coinAmount, fromUserId}
		addQueryValues         = []any{coinAmount, toUserId}
		transactionQueryValues = []any{fromUserId, toUserId, coinAmount}
	)

	//TODO: think about isolation level
	txOptions := pgx.TxOptions{IsoLevel: pgx.Serializable}
	tx, err := r.dbPool.BeginTx(ctx, txOptions)
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	res, err := tx.Exec(ctx, subtractQuery, subtractQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	res, err = tx.Exec(ctx, addQuery, addQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	_, err = tx.Exec(ctx, insertTransactionQuery, transactionQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	return nil
}
