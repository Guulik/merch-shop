package service

import "context"

type CoinTransfer interface {
	TransferCoins(
		ctx context.Context,
		fromUserId int,
		toUserId int,
		coinAmount int,
	) error
}

func (s *Service) SendCoins(ctx context.Context, fromUserId int, toUserId int, coinAmount int) error {
	//TODO: implement me
	return nil
}
