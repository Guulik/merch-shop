package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
	"net/http"
)

type UserProvider interface {
	GetCoins(
		ctx context.Context,
		userId int,
	) (int, error)
	GetCoinsAndInventory(
		ctx context.Context,
		userId int,
	) (*int, map[string]int, error)
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
		coinsPtr     *int
		inventoryMap map[string]int
		inventory    []model.Item
		coinHistory  model.CoinHistory
		err          error
	)

	coinsPtr, inventoryMap, err = s.userProvider.GetCoinsAndInventory(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserInfo{}, echo.NewHTTPError(http.StatusBadRequest, consts.UserNotFound)
		}
		return &model.UserInfo{}, echo.NewHTTPError(http.StatusInternalServerError, consts.InternalServerError)
	}
	if coinsPtr != nil {
		coins = *coinsPtr
	}
	logger.WithLogCoinBalance(ctx, coins)

	inventory = convertInventory(inventoryMap)

	coinHistory, err = s.userProvider.GetCoinHistory(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserInfo{}, echo.NewHTTPError(http.StatusBadRequest, consts.UserNotFound)
		}
		return &model.UserInfo{}, echo.NewHTTPError(http.StatusInternalServerError, consts.InternalServerError)
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
