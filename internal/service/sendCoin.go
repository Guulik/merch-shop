package service

import "context"

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
		//TODO: handle error
	}
	if currentCoins < coinAmount {
		//TODO: log and return 400 (consts.NotEnoughMoney)
	}

	err = s.coinTransfer.TransferCoins(ctx, fromUserId, toUserId, coinAmount)
	if err != nil {
		//TODO: handle error
	}

	return nil
}
