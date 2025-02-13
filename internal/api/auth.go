package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
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
	logger.WithLogUsername(ctx, req.Username)

	slog.Debug("api authorize")
	token, err = a.service.Authorize(ctx, req.Username, req.Password)
	if err != nil {
		var httpErr *wrapper.HTTPError
		if errors.As(err, &httpErr) {
			slog.ErrorContext(logger.ErrorCtx(ctx, httpErr.Err), "Error: "+httpErr.Err.Error())
			return echo.NewHTTPError(httpErr.Status, httpErr.Msg)
		}
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
		return err
	}

	return e.JSON(http.StatusOK, AuthResponse{Token: token})
}
