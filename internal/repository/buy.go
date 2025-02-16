package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"

	"merch/internal/util/logger"
)

func (r *Repo) PayForItem(ctx context.Context, userId int, item string, itemCost int) error {
	var (
		coinQuery = `
		UPDATE users 
		SET coins = coins - $1
		WHERE id = $2 AND coins >= $1
`

		itemQuery = `
		INSERT INTO inventory (user_id, item, quantity) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, item) 
		DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity;
`

		coinQueryValues = []any{itemCost, userId}
		itemQueryValues = []any{userId, item, 1}
	)

	txOptions := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	tx, err := r.dbPool.BeginTx(ctx, txOptions)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	res, err := tx.Exec(ctx, coinQuery, coinQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	res, err = tx.Exec(ctx, itemQuery, itemQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	return nil
}
