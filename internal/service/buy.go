package service

import (
	"context"
	"database/sql"
	"errors"
	"merch/internal/domain"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
)

func (s *Service) BuyItem(ctx context.Context, userId int, item string) error {

	var (
		currentCoins int
		err          error
	)

	//Уже провалидировано на api уровне
	itemCost := domain.Shop[item]

	currentCoins, err = s.userProvider.GetCoins(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: return 400
		}
		//TODO: return 500
	}
	logger.WithLogCoinBalance(ctx, currentCoins)
	if currentCoins < itemCost {
		err = errors.New(consts.NotEnoughMoney)
		return logger.WrapError(ctx, err)
	}

	err = s.coinTransfer.PayForItem(ctx, userId, item, itemCost)
	if err != nil {
		//Возможно избыточно, ведь мы проверили пользователя в GetCoins, но я решил оставить. Это небольшая нагрузка
		if errors.Is(err, sql.ErrNoRows) {
			//TODO: return 400
		}
		//TODO: return 500
	}

	return nil
}
