package auth

import (
	"github.com/dchaykin/go-modules/log"
	"github.com/golang-jwt/jwt/v4"
)

type UserIdentity interface {
	Partner() string
}

type userIdentity struct {
	token  *jwt.Token
	claims jwt.MapClaims
}

func (j userIdentity) Partner() string {
	claim, ok := j.claims["partner"]
	if !ok {
		log.Warn("Claim 'partner' not found")
		return ""
	}
	return claim.(string)
}

func GetUserIdentity(authorization, secret string) (UserIdentity, error) {
	t, claims, err := parseToken(authorization, secret)
	if err != nil {
		return nil, err
	}
	return &userIdentity{
		token:  t,
		claims: claims,
	}, nil
}

func CreateAuthorizationToken(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
