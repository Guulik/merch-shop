package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
	"strings"
)

func (a *Api) InfoHandler(e echo.Context) error {
	ctx := e.Request().Context()
	var (
		userInfo    *model.UserInfo
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	ctx = logger.WithLogUserID(ctx, tokenUserId)

	userInfo, err = a.service.GetUserInfo(ctx, tokenUserId)
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

	return e.JSON(http.StatusOK, userInfo)
}
