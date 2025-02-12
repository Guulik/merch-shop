package service

import (
	"context"
	"merch/internal/domain"
)

func (s *Service) BuyItem(ctx context.Context, userId int, item string) error {

	var (
		currentCoins int
		err          error
	)

	itemCost, ok := domain.Shop[item]
	if !ok {
		//TODO: log and return 500
	}
	currentCoins, err = s.userProvider.GetCoins(ctx, userId)
	if err != nil {
		//TODO: handle error
	}
	if currentCoins < itemCost {
		//TODO: log and return 400 (consts.NotEnoughMoney)
	}

	err = s.coinTransfer.PayForItem(ctx, userId, item, itemCost)
	if err != nil {
		//TODO: handle error
	}

	return nil
}
