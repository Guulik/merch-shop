package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"log/slog"
	"merch/internal/domain/consts"
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
	) (int, error)
	SaveUser(
		ctx context.Context,
		username string,
		password string,
	) (int, error)
}

func (s *Service) Authorize(ctx context.Context, username string, password string) (string, error) {

	var (
		hashedPassword string
		userId         int
		newUserId      int

		token string
		err   error
	)
	hashedPassword, err = hasher.HashPassword(password)
	if err != nil {
		return "", logger.WrapError(ctx, err)
	}

	userId, err = s.authorizer.CheckUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("new user")
			newUserId, err = s.authorizer.SaveUser(ctx, username, hashedPassword)
			if err != nil {
				return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
			}
		} else {
			return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
		}
		token, err = s.generateJWT(newUserId)
		if err != nil {
			return "", logger.WrapError(ctx, err)
		}
		return token, nil
	}

	if err = hasher.ComparePassword(hashedPassword, password); err != nil {
		return "", wrapper.WrapHTTPError(err, http.StatusUnauthorized, consts.WrongPassword)
	}

	token, err = s.generateJWT(userId)
	if err != nil {
		return "", logger.WrapError(ctx, err)
	}
	return token, nil
}

func (s *Service) generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.cfg.TokenTTL).Unix(),
	}

	//TODO: choose signing method
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
