package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4"

	"merch/internal/domain"
	"merch/internal/domain/consts"
	"merch/internal/util/logger"
	"merch/internal/util/wrapper"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	ctx = logger.WithLogCoinBalance(ctx, currentCoins)
	if currentCoins < itemCost {
		err = errors.New(consts.NotEnoughMoney)
		return logger.WrapError(ctx, err)
	}

	err = s.coinTransfer.PayForItem(ctx, userId, item, itemCost)
	if err != nil {
		//Возможно, это избыточно, ведь мы проверили пользователя в GetCoins, но накладных расходов почти не создает.
		if errors.Is(err, pgx.ErrNoRows) {
			return wrapper.WrapHTTPError(err, http.StatusBadRequest, consts.UserNotFound)
		}
		return wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}

	return nil
}
