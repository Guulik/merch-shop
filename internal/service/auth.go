package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"merch/internal/domain/model"
	"time"
)

type Authorizer interface {
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (model.UserAuth, error)
	SaveUser(
		ctx context.Context,
		username string,
		password string,
	) error
}

func (s *Service) Authorize(ctx context.Context, username string, password string) (string, error) {

	var (
		token string
		user  model.UserAuth
		err   error
	)
	user, err = s.authorizer.GetUserByUsername(ctx, username)
	if err != nil {
		//TODO: check if error "user not found", create new User, save password as hash
	}
	//TODO: check password
	/*	if err = utils.ComparePassword(password, user.PasswordHash); err != nil {
		//TODO: return 401
	}*/

	token, err = s.generateJWT(user.Id)
	if err != nil {
		//TODO: return 500
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
