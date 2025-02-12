package repository

import (
	"context"
	"merch/internal/domain/model"
)

func (r *Repo) GetCoinsAndInventory(ctx context.Context, userId int) (*int, map[string]int, error) {

	var (
		query = `
    SELECT u.coins as coins, i.item, i.quantity
    FROM users u
    JOIN inventory i on u.id = i.user_id
    WHERE u.id = $1
`
		values = []any{userId}

		coins     int
		inventory map[string]int
	)

	rows, err := r.dbPool.Query(query, values...)
	if err != nil {
		//TODO: handle error
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		var quantity int

		if err = rows.Scan(&coins, &item, &quantity); err != nil {
			//TODO: handle error
		}

		if quantity != 0 {
			inventory[item] = quantity
		}
	}

	return &coins, inventory, nil
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
	err := r.dbPool.Get(&transactions, query, values...)
	if err != nil {
		//TODO: handle error
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
