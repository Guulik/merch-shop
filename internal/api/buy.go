package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"merch/internal/domain"
	"merch/internal/domain/consts"
	"net/http"
)

type BuyRequest struct {
	Item string `query:"item"`
}

func (a *Api) BuyHandler(e echo.Context) error {
	const op = "Api.BuyHandler"

	var (
		req         BuyRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	ctx := e.Request().Context()

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}

	err = validateItem(req.Item)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = a.service.BuyItem(ctx, tokenUserId, req.Item)
	if err != nil {
		//TODO: return wrapped error
	}

	return e.NoContent(http.StatusOK)
}

func validateItem(item string) error {
	for _, i := range domain.Items {
		if item == i {
			return nil
		}
	}
	return errors.New(consts.WrongItem)
}
