package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/lib/secret"
	"net/http"
	"strings"
)

const (
	InvalidTokenFormat = "invalid token format"
	MissingToken       = "missing token"
	InvalidToken       = "invalid token"
	InvalidTokenClaims = "invalid token claims"
	UserIdNotFound     = "user_id not found in token"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			slog.Debug("no auth header")
			return echo.NewHTTPError(http.StatusUnauthorized, MissingToken)
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Debug("no token bearer")
			return echo.NewHTTPError(http.StatusUnauthorized, InvalidTokenFormat)
		}
		tokenString := parts[1]

		secretKey, err := secret.FetchSecretKey()
		if err != nil {
			slog.Error("failed to fetch secret key: ", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, consts.InternalServerError)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			slog.Debug("invalid token", slog.String("token:", token.Raw))
			return echo.NewHTTPError(http.StatusUnauthorized, InvalidToken)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			slog.Debug("invalid token claims", slog.Any("claims:", token.Claims))
			return echo.NewHTTPError(http.StatusUnauthorized, InvalidTokenClaims)
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			slog.Debug("user claims not found", slog.Any("claims:", claims))
			return echo.NewHTTPError(http.StatusUnauthorized, UserIdNotFound)
		}

		c.Set("user_id", int(userID))

		return next(c)
	}
}
