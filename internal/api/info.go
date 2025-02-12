package api

import (
	"github.com/labstack/echo/v4"
	"merch/internal/domain/model"
	"net/http"
)

func (a *Api) InfoHandler(e echo.Context) error {
	const op = "Api.InfoHandler"

	var (
		userInfo    *model.UserInfo
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)

	ctx := e.Request().Context()

	userInfo, err = a.service.GetUserInfo(ctx, tokenUserId)
	if err != nil {
		//TODO: return wrapped error
	}

	return e.JSON(http.StatusOK, userInfo)
}
