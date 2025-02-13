package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
	"net/http"
)

type BuyRequest struct {
	Item string `query:"item"`
}

func (a *Api) BuyHandler(e echo.Context) error {
	ctx := e.Request().Context()

	var (
		req         BuyRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	logger.WithLogUserID(ctx, tokenUserId)

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}

	err = validateItem(req.Item)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	logger.WithLogItem(ctx, req.Item)

	err = a.service.BuyItem(ctx, tokenUserId, req.Item)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error:"+err.Error())
		//TODO: return wrapped error
	}

	return e.NoContent(http.StatusOK)
}

func validateItem(item string) error {
	_, ok := domain.Shop[item]
	if !ok {
		return errors.New(consts.WrongItem)
	}
	return nil
}
