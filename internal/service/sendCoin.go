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

func (s *Service) SendCoins(ctx context.Context, fromUserId int, toUsername string, coinAmount int) error {

	var (
		toUser       *model.UserAuth
		currentCoins int
		err          error
	)
	toUser, err = s.authorizer.CheckUserByUsername(ctx, toUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if err != nil {
				return wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.ToUserNotFound)
			}
		} else {
			return wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
		}
	}

	currentCoins, err = s.userProvider.GetCoins(ctx, fromUserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}
	ctx = logger.WithLogCoinBalance(ctx, currentCoins)
	if currentCoins < coinAmount {
		err = errors.New(consts.NotEnoughMoney)
		return wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.NotEnoughMoney)
	}

	err = s.coinTransfer.TransferCoins(ctx, fromUserId, toUser.Id, coinAmount)
	if err != nil {
		return wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	return nil
}
