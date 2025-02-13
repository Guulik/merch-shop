package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
)

type SendCoinRequest struct {
	ToUser int `json:"toUser"`
	Amount int `json:"amount"`
}

func (a *Api) SendCoinHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var (
		req         SendCoinRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	logger.WithLogUserID(ctx, tokenUserId)

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}
	logger.WithLogToUser(ctx, req.ToUser)
	logger.WithLogSendAmount(ctx, req.Amount)

	err = a.service.SendCoins(ctx, tokenUserId, req.ToUser, req.Amount)
	if err != nil {
		var httpErr *wrapper.HTTPError
		if errors.As(err, &httpErr) {
			slog.ErrorContext(logger.ErrorCtx(ctx, httpErr.Err), "Error: "+httpErr.Err.Error())
			return echo.NewHTTPError(httpErr.Status, httpErr.Msg)
		}
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
	}

	return e.NoContent(http.StatusOK)
}
