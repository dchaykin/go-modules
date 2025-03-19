package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func parseToken(tokenString, secret string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return token, claims, nil
	}
	return nil, nil, fmt.Errorf("invalid token: %t, %t", ok, token.Valid)
}
