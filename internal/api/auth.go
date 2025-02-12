package api

import (
	"github.com/labstack/echo/v4"
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
	const op = "Api.AuthHandler"
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
	token, err = a.service.Authorize(ctx, req.Username, req.Password)

	return e.JSON(http.StatusOK, AuthResponse{Token: token})
}
