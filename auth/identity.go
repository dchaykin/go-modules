package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dchaykin/go-modules/log"
	"github.com/golang-jwt/jwt/v4"
)

type UserIdentity interface {
	Partner() string
	Role(name string) string
	FirstName() string
	SurName() string
	Email() string
	Username() string
}

type userToken struct {
	Claims jwt.MapClaims `json:"claims"`
}

func (j userToken) Partner() string {
	claim, ok := j.Claims["partner"]
	if !ok {
		log.Warn("Claim 'partner' not found")
		return ""
	}
	return claim.(string)
}

func (j userToken) Role(name string) string {
	var claim interface{}
	var ok bool
	switch name {
	case "inquiry":
		claim, ok = j.Claims["roleInquiry"]
	default:
		claim, ok = j.Claims["role"]
	}
	if !ok {
		log.Warn("Claim 'role' for %s not found", name)
		return ""
	}
	return claim.(string)
}

func (j userToken) FirstName() string {
	claim, ok := j.Claims["firstName"]
	if !ok {
		log.Warn("Claim 'firstName' not found")
		return ""
	}
	return claim.(string)
}

func (j userToken) SurName() string {
	claim, ok := j.Claims["surName"]
	if !ok {
		log.Warn("Claim 'surName' not found")
		return ""
	}
	return claim.(string)
}

func (j userToken) Email() string {
	claim, ok := j.Claims["eMail"]
	if !ok {
		log.Warn("Claim 'eMail' not found")
		return ""
	}
	return claim.(string)
}

func (j userToken) Username() string {
	claim, ok := j.Claims["userName"]
	if !ok {
		log.Warn("Claim 'userName' not found")
		return ""
	}
	return claim.(string)
}

func GetUserIdentityFromRequest(r http.Request) (UserIdentity, error) {
	userInfo := r.Header.Get("X-User-Info")
	if userInfo == "" {
		return nil, fmt.Errorf("no user info in the request found")
	}
	ui := userToken{}
	err := json.Unmarshal([]byte(userInfo), &ui)
	return ui, err
}

func GetUserIdentity(authorization, secret string) (UserIdentity, error) {
	claims, err := parseToken(authorization, secret)
	if err != nil {
		return nil, err
	}
	return &userToken{
		Claims: claims,
	}, nil
}

func CreateAuthorizationToken(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
