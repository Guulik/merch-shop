package jwtManager

import (
	"errors"
	"os"
)

func FetchSecretKey() ([]byte, error) {
	const key = "JWT_SECRET"

	if secret := os.Getenv(key); secret != "" {
		return []byte(secret), nil
	}

	return nil, errors.New("JWT_SECRET is not set")
}
