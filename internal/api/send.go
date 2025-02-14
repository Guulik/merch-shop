package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
	"strings"
)

type SendCoinRequest struct {
	ToUsername string `json:"toUser"`
	Amount     int    `json:"amount"`
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
	logger.WithLogToUser(ctx, req.ToUsername)
	logger.WithLogSendAmount(ctx, req.Amount)

	err = a.service.SendCoins(ctx, tokenUserId, req.ToUsername, req.Amount)
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
	}

	return e.NoContent(http.StatusOK)
}
