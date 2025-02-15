package repository

import (
	"context"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
)

func (r *Repo) GetCoins(ctx context.Context, userId int) (int, error) {
	var (
		query = `
    SELECT coins
    FROM users
    WHERE id = $1
`
		values = []any{userId}

		coins int
	)

	err := r.dbPool.QueryRow(ctx, query, values...).Scan(&coins)
	if err != nil {
		return -1, logger.WrapError(ctx, err)
	}

	return coins, nil
}

func (r *Repo) GetInventory(ctx context.Context, userId int) (map[string]int, error) {
	var (
		query = `
    SELECT item, quantity
    FROM inventory
    WHERE user_id = $1
`
		values = []any{userId}

		inventory = make(map[string]int)
	)

	rows, err := r.dbPool.Query(ctx, query, values...)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		var quantity int

		if err = rows.Scan(&item, &quantity); err != nil {
			return nil, logger.WrapError(ctx, err)
		}

		if quantity != 0 {
			inventory[item] = quantity
		}
	}

	return inventory, nil
}

func (r *Repo) GetCoinHistory(ctx context.Context, userId int) (model.CoinHistory, error) {
	var (
		query = `
	SELECT 
    t.from_user AS fromUserId,
    from_user.username AS fromUsername,
    t.to_user AS toUserId,
    to_user.username AS toUsername,
    t.amount
	FROM transactions t
	JOIN users from_user ON t.from_user = from_user.id
	JOIN users to_user ON t.to_user = to_user.id
	WHERE t.from_user = $1 OR t.to_user = $1;
`
		values       = []any{userId}
		transactions []model.Transaction
	)

	var (
		coinHistory model.CoinHistory
		received    []model.Received
		sent        []model.Sent
	)

	rows, err := r.dbPool.Query(ctx, query, values...)
	if err != nil {
		return model.CoinHistory{}, logger.WrapError(ctx, err)
	}
	defer rows.Close()

	for rows.Next() {
		var utx model.Transaction
		if err = rows.Scan(
			&utx.FromUserId,
			&utx.FromUsername,
			&utx.ToUserId,
			&utx.ToUsername,
			&utx.Amount,
		); err != nil {
			return model.CoinHistory{}, logger.WrapError(ctx, err)
		}
		transactions = append(transactions, utx)
	}

	if err = rows.Err(); err != nil {
		return model.CoinHistory{}, logger.WrapError(ctx, err)
	}

	if err != nil {
		return model.CoinHistory{}, logger.WrapError(ctx, err)
	}
	for _, t := range transactions {
		if t.ToUserId == userId {
			received = append(received, model.Received{
				FromUser: t.FromUsername,
				Amount:   t.Amount,
			})
		}
		if t.FromUserId == userId {
			sent = append(sent, model.Sent{
				ToUser: t.ToUsername,
				Amount: t.Amount,
			})
		}
	}

	coinHistory.Received = received
	coinHistory.Sent = sent

	return coinHistory, nil
}
