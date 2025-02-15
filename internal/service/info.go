package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
)

type UserProvider interface {
	GetCoins(
		ctx context.Context,
		userId int,
	) (int, error)
	GetInventory(
		ctx context.Context,
		userId int,
	) (map[string]int, error)
	GetCoinHistory(
		ctx context.Context,
		userId int,
	) (model.CoinHistory, error)
}

// GetUserInfo использует 2 запроса вместо 3.
// В одном запросе получаем сразу и монетки и инвентарь, чтобы лишний раз не ходить в базу за монетками.
// Во втором запросе все транзакции, относящиеся к данному пользователю, но разбитые на sent и received
func (s *Service) GetUserInfo(ctx context.Context, userId int) (*model.UserInfo, error) {
	var (
		coins        int
		inventoryMap map[string]int
		inventory    []model.Item
		coinHistory  model.CoinHistory
		err          error
	)

	coins, err = s.userProvider.GetCoins(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	inventoryMap, err = s.userProvider.GetInventory(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	ctx = logger.WithLogCoinBalance(ctx, coins)

	inventory = convertInventory(inventoryMap)

	coinHistory, err = s.userProvider.GetCoinHistory(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return &model.UserInfo{}, wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	userInfo := &model.UserInfo{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	return userInfo, nil
}

func convertInventory(inventoryMap map[string]int) []model.Item {
	inventory := make([]model.Item, 0, len(inventoryMap))
	for item, quantity := range inventoryMap {
		inventory = append(inventory, model.Item{
			Type:     item,
			Quantity: quantity,
		})
	}
	return inventory
}
