package jwtManager

import "github.com/golang-jwt/jwt/v5"

func ParseToken(token string, secretKey []byte) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
