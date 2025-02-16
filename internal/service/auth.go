package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"

	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/util/hasher"
	"merch/internal/util/jwtManager"
	"merch/internal/util/wrapper"
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
				return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
			}
			slog.Debug("new user", slog.String("hashed password", hashedPassword))
			newUserId, err = s.authorizer.SaveUser(ctx, username, hashedPassword)
			if err != nil {
				return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
			}
		} else {
			return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
		}
		token, err = jwtManager.GenerateJWT(newUserId, s.cfg.TokenTTL)
		if err != nil {
			return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
		}
		return token, nil
	}

	slog.Debug("password from db: " + user.PasswordDb)
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordDb), []byte(password)); err != nil {
		return "", wrapper.WrapHTTPError(err, http.StatusUnauthorized, consts.WrongPassword)
	}

	token, err = jwtManager.GenerateJWT(user.Id, s.cfg.TokenTTL)
	if err != nil {
		return "", wrapper.WrapHTTPError(err, http.StatusInternalServerError, consts.InternalServerError)
	}
	return token, nil
}
