package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/hasher"
	"merch/internal/lib/logger"
	"merch/internal/lib/secret"
	"merch/internal/lib/wrapper"
	"net/http"
	"time"
)

type Authorizer interface {
	CheckUserByUsername(
		ctx context.Context,
		username string,
	) (*model.UserAuth, error)
	SaveUser(
		ctx context.Context,
		username string,
		password string,
	) (int, error)
}

func (s *Service) Authorize(ctx context.Context, username string, password string) (string, error) {

	var (
		user      *model.UserAuth
		newUserId int

		token string
		err   error
	)

	user, err = s.authorizer.CheckUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hashedPassword, err := hasher.HashPassword(password)
			if err != nil {
				return "", logger.WrapError(ctx, err)
			}
			slog.Debug("new user", slog.String("hashed password", hashedPassword))
			newUserId, err = s.authorizer.SaveUser(ctx, username, hashedPassword)
			if err != nil {
				return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
			}
		} else {
			return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
		}
		token, err = s.generateJWT(newUserId, s.tokenTTL)
		if err != nil {
			return "", logger.WrapError(ctx, err)
		}
		return token, nil
	}

	slog.Debug("password from db: " + user.PasswordDb)
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordDb), []byte(password)); err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, consts.WrongPassword)
	}

	token, err = s.generateJWT(user.Id, s.tokenTTL)
	if err != nil {
		return "", logger.WrapError(ctx, err)
	}
	return token, nil
}

func (s *Service) generateJWT(userID int, TTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(TTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := secret.FetchSecretKey()
	if err != nil {
		return "", err
	}

	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}
