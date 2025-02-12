package service

import "context"

type Buyer interface {
	PayForItem(
		ctx context.Context,
		userId int,
		item string,
		itemCost int,
	) error
}

func (s *Service) BuyItem(ctx context.Context, userId int, item string) error {
	//TODO: implement me
	return nil
}
