package api

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/lib/logger"
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
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error:"+err.Error())
		//TODO: return wrapped error
	}

	return e.NoContent(http.StatusOK)
}
