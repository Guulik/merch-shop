package repository

import (
	"context"
	"database/sql"
	"merch/internal/lib/logger"
)

func (r *Repo) PayForItem(ctx context.Context, userId int, item string, itemCost int) error {
	//TODO: wrap sql with squirrel
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

	//TODO: replace to BeginTx. Provide context and choose iso level
	tx, err := r.dbPool.Begin()
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, coinQuery, coinQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	res, err = tx.ExecContext(ctx, itemQuery, itemQueryValues...)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
		return logger.WrapError(ctx, err)
	}

	err = tx.Commit()
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	return nil
}
