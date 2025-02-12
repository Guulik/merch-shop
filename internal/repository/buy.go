package repository

import (
	"context"
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

	tx, err := r.dbPool.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(coinQuery, coinQueryValues...)
	if err != nil {
		//TODO: log error
		return err
	}

	_, err = tx.Exec(itemQuery, itemQueryValues...)
	if err != nil {
		//TODO: log error
		return err
	}

	return tx.Commit()
}
