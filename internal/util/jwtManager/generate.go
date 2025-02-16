package jwtManager

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int, TTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(TTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := FetchSecretKey()
	if err != nil {
		return "", err
	}

	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}
