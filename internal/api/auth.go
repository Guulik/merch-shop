package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/lib/logger"
	"merch/internal/lib/wrapper"
	"net/http"
	"strings"
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
