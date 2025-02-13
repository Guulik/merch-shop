package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"log/slog"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/hasher"
	"merch/internal/lib/logger"
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
		hashedPassword string
		user           *model.UserAuth
		newUserId      int

		token string
		err   error
	)
	hashedPassword, err = hasher.HashPassword(password)
	if err != nil {
		return "", logger.WrapError(ctx, err)
	}

	user, err = s.authorizer.CheckUserByUsername(ctx, username)
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

	if err = hasher.ComparePassword(hashedPassword, user.Password); err != nil {
		return "", wrapper.WrapHTTPError(err, http.StatusUnauthorized, consts.WrongPassword)
	}

	token, err = s.generateJWT(user.Id)
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

	key, err := fetchSecretKey()
	if err != nil {
		return "", err
	}

	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}

func fetchSecretKey() ([]byte, error) {
	//TODO: implement me!
	return []byte("kkkk"), nil
}
