package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
	"strings"
)

type BuyRequest struct {
	Item string `param:"item" validate:"required"`
}

func (a *Api) BuyHandler(e echo.Context) error {
	ctx := e.Request().Context()

	var (
		req         BuyRequest
		tokenUserId int
		err         error
	)
	tokenUserId = e.Get("user_id").(int)
	ctx = logger.WithLogUserID(ctx, tokenUserId)

	if err = e.Bind(&req); err != nil {
		// always returns wrapped 400
		return err
	}
	slog.Debug("item: " + req.Item)
	err = validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = validateItem(req.Item)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}
	ctx = logger.WithLogItem(ctx, req.Item)

	err = a.service.BuyItem(ctx, tokenUserId, req.Item)
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

	return e.NoContent(http.StatusOK)
}

func validateItem(item string) error {
	_, ok := domain.Shop[item]
	if !ok {
		return errors.New(consts.WrongItem)
	}
	return nil
}
