package repository

import (
	"context"
	"database/sql"
)

func (r *Repo) TransferCoins(ctx context.Context, fromUserId int, toUserId int, coinAmount int) error {

	var (
		subtractQuery = `
	UPDATE users 
	SET coins = coins - $1
	WHERE id = $2 AND coins >= $1
	RETURNING coins
`
		addQuery = `
	UPDATE users 
	SET coins = coins + $1
	WHERE id = $2
	RETURNING coins
`

		insertTransactionQuery = `
	INSERT INTO transactions (from_user, to_user, amount) 
	VALUES ($1, $2, $3)
`

		subtractQueryValues    = []any{coinAmount, fromUserId}
		addQueryValues         = []any{coinAmount, toUserId}
		transactionQueryValues = []any{fromUserId, toUserId, coinAmount}
	)

	txOptions := sql.TxOptions{Isolation: sql.LevelSerializable}
	tx, err := r.dbPool.BeginTx(ctx, &txOptions)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(subtractQuery, subtractQueryValues...)
	if err != nil {
		return err
	}

	_, err = tx.Exec(addQuery, addQueryValues...)
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertTransactionQuery, transactionQueryValues...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
