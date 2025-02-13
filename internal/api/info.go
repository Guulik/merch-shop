package api

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/model"
	"merch/internal/lib/logger"
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
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error:"+err.Error())
		//TODO: return wrapped error
	}

	return e.JSON(http.StatusOK, userInfo)
}
