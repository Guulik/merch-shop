package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"merch/internal/domain/consts"
	"merch/internal/util/logger"
	"merch/internal/util/wrapper"
)

type SendCoinRequest struct {
	ToUsername string `json:"toUser" validate:"required"`
	Amount     int    `json:"amount" validate:"required,gt=0"`
}

type CoinSender interface {
	SendCoins(ctx context.Context, fromUserId int, toUsername string, coinAmount int) error
}

func (a *Api) SendCoinHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var (
		req         SendCoinRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	ctx = logger.WithLogUserID(ctx, tokenUserId)

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}
	err = validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx = logger.WithLogToUser(ctx, req.ToUsername)
	ctx = logger.WithLogSendAmount(ctx, req.Amount)

	err = a.coinSender.SendCoins(ctx, tokenUserId, req.ToUsername, req.Amount)
	if err != nil {
		var httpErr *wrapper.HTTPError
		if errors.As(err, &httpErr) {
			// Не уверен. О сравнении надо ещё подумать
			if strings.Contains(err.Error(), consts.InternalServerError) {
				slog.ErrorContext(logger.ErrorCtx(ctx, httpErr.InternalErr), "Error: "+httpErr.InternalErr.Error())
			} else {
				slog.WarnContext(logger.ErrorCtx(ctx, httpErr.InternalErr), "Error: "+httpErr.InternalErr.Error())
			}
			return echo.NewHTTPError(httpErr.Status, httpErr.Msg)
		}
		slog.WarnContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
		return err
	}

	return e.NoContent(http.StatusOK)
}
