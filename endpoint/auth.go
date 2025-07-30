package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/httpcomm"
	"github.com/dchaykin/go-modules/log"
)

type token struct {
	Token string `json:"token"`
}

func CreateUserIdentityByCredentials(username, password, secret string) (auth.UserIdentity, error) {
	payload := fmt.Appendf(nil, `{"username":"%s","password":"%s"}`, username, password)
	endpoint := fmt.Sprintf("https://%s/app-config/auth", os.Getenv("MYHOST"))
	resp := httpcomm.Post(endpoint, nil, nil, payload)
	if resp.StatusCode != http.StatusOK {
		return nil, log.WrapError(resp.GetError())
	}

	t := token{}
	err := json.Unmarshal(resp.Answer, &t)
	if err != nil {
		return nil, log.WrapError(err)
	}

	return auth.GetUserIdentity(t.Token, secret)
}
