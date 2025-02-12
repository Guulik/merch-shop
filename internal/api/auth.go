package api

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"time"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (a *Api) AuthHandler(c echo.Context) error {
	const op = "Api.AuthHandler"

	var (
		req   AuthRequest
		token string
		err   error
	)
	if err := c.Bind(&req); err != nil {
		//TODO: return 400
	}

	//TODO: check username
	//user, err := a.service.GetUserByUsername(req.Username)
	if err != nil {
		//TODO: return 401
	}

	//TODO: check password
	/*	if err = utils.ComparePassword(req.Password, user.PasswordHash); err != nil {
		//TODO: return 401
	}*/

	//TODO: generate JWT based on userId
	/*	token, err = a.generateJWT(user.ID)
		if err != nil {
			//TODO: return 400
		}*/

	// Todo: Return 200 + token
	return c.JSON(200, AuthResponse{Token: token})
}

func (a *Api) generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(a.cfg.TokenTTL).Unix(),
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
