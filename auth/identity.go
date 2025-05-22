package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	IsAdmin() bool
	IsDeveloper() bool
	Set(req *http.Request) error
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
	claim, ok := j.Claims[name]
	if !ok {
		log.Warn("Role '%s' is not claimed", name)
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

func (j userToken) IsAdmin() bool {
	claim, ok := j.Claims["admin"]
	if !ok {
		return false
	}
	return claim.(bool)
}

func (j userToken) IsDeveloper() bool {
	if j.Username() == "dchaykin" { // TODO
		return true
	}
	claim, ok := j.Claims["developer"]
	if !ok {
		return false
	}
	return claim.(bool)
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

func (j userToken) Set(req *http.Request) error {
	authorization, err := CreateAuthorizationToken(j.Claims, os.Getenv("AUTH_SECRET"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(authorization))
	return nil
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
