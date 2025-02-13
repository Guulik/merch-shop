package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
)

func (a *Api) InfoHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var (
		userInfo    *model.UserInfo
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	logger.WithLogUserID(ctx, tokenUserId)

	userInfo, err = a.service.GetUserInfo(ctx, tokenUserId)
	if err != nil {
		var httpErr *wrapper.HTTPError
		if errors.As(err, &httpErr) {
			slog.ErrorContext(logger.ErrorCtx(ctx, httpErr.Err), "Error: "+httpErr.Err.Error())
			return echo.NewHTTPError(httpErr.Status, httpErr.Msg)
		}
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
	}

	return e.JSON(http.StatusOK, userInfo)
}
