package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4"

	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/util/logger"
	"merch/internal/util/wrapper"
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
