package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type SendCoinRequest struct {
	ToUser int `json:"toUser"`
	Amount int `json:"amount"`
}

func (a *Api) SendCoinHandler(e echo.Context) error {
	const op = "Api.SendCoinHandler"

	var (
		req         SendCoinRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	ctx := e.Request().Context()

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}
	err = a.service.SendCoins(ctx, tokenUserId, req.ToUser, req.Amount)
	if err != nil {
		//TODO: return wrapped error
	}

	return e.NoContent(http.StatusOK)
}
