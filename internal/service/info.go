package service

import (
	"context"
	"merch/internal/domain/model"
)

// UserProvider содержит 2 запроса вместо 3 для оптимизации.
// В одном запросе получаем сразу и монетки и инвентарь.
// Во втором запросе все транзакции, относящиеся к данному пользователю, но разбитые на sent и received
type UserProvider interface {
	GetCoinsAndInventory(ctx context.Context, userId int) (*int, map[string]int, error)
	GetCoinHistory(
		ctx context.Context,
		userId int,
	) (model.CoinHistory, error)
}

func (s *Service) GetUserInfo(ctx context.Context, userId int) (*model.UserInfo, error) {
	var (
		coins        int
		coinsPtr     *int
		inventoryMap map[string]int
		inventory    []model.Item
		coinHistory  model.CoinHistory
		err          error
	)

	coinsPtr, inventoryMap, err = s.userProvider.GetCoinsAndInventory(ctx, userId)
	if err != nil {
		//TODO: handle error
	}
	if coinsPtr != nil {
		coins = *coinsPtr
	}
	//TODO: convert map[string]int to []Item
	inventory, err = convertInventory(inventoryMap)

	coinHistory, err = s.userProvider.GetCoinHistory(ctx, userId)
	if err != nil {
		//TODO: handle error
	}

	userInfo := &model.UserInfo{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return userInfo, nil
}

func convertInventory(inventoryMap map[string]int) ([]model.Item, error) {
	//TODO: implement me
	return nil, nil
}
