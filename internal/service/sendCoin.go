package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
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

func (s *Service) SendCoins(ctx context.Context, fromUserId int, toUserId int, coinAmount int) error {

	var (
		currentCoins int
		err          error
	)

	currentCoins, err = s.userProvider.GetCoins(ctx, fromUserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, consts.UserNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, consts.InternalServerError)
	}
	logger.WithLogCoinBalance(ctx, currentCoins)
	if currentCoins < coinAmount {
		err = errors.New(consts.NotEnoughMoney)
		return logger.WrapError(ctx, err)
	}

	err = s.coinTransfer.TransferCoins(ctx, fromUserId, toUserId, coinAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, consts.UserNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, consts.InternalServerError)
	}

	return nil
}
