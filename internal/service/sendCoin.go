package service

import (
	"context"
	"database/sql"
	"errors"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
)

type CoinTransfer interface {
	TransferCoins(
		ctx context.Context,
		fromUserId int,
		toUserId int,
		coinAmount int,
	) error
	PayForItem(
		ctx context.Context,
		userId int,
		item string,
		itemCost int,
	) error
}

func (s *Service) SendCoins(ctx context.Context, fromUserId int, toUserId int, coinAmount int) error {

	var (
		currentCoins int
		err          error
	)

	currentCoins, err = s.userProvider.GetCoins(ctx, fromUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: return 400
		}
		//TODO: return 500
	}
	logger.WithLogCoinBalance(ctx, currentCoins)
	if currentCoins < coinAmount {
		err = errors.New(consts.NotEnoughMoney)
		return logger.WrapError(ctx, err)
	}

	err = s.coinTransfer.TransferCoins(ctx, fromUserId, toUserId, coinAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: return 400
		}
		//TODO: return 500
	}

	return nil
}
