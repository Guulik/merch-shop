package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"merch/internal/domain/consts"
	"merch/internal/util/logger"
	"merch/internal/util/wrapper"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthorizerService interface {
	Authorize(ctx context.Context, username string, password string) (string, error)
}

func (a *Api) AuthHandler(e echo.Context) error {
	ctx := e.Request().Context()

	var (
		req   AuthRequest
		token string
		err   error
	)
	if err = e.Bind(&req); err != nil {
		// always returns 400
		return err
	}
	ctx = logger.WithLogUsername(ctx, req.Username)
	err = validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	slog.Debug("api authorize")
	token, err = a.authorizer.Authorize(ctx, req.Username, req.Password)
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

	return e.JSON(http.StatusOK, AuthResponse{Token: token})
}
