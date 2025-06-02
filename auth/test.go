package auth

import (
	"fmt"
	"net/http"
)

type TestUser struct {
	Claims map[string]any
	tenant string
}

func (u TestUser) Partner() string {
	return u.Claims["partner"].(string)
}

func (u TestUser) Tenant() string {
	return u.tenant
}

func (u TestUser) RoleBySubject(subject string) string {
	roles := u.Claims["roles"].(map[string]any)
	return fmt.Sprintf("%v", roles[subject])
}

func (u TestUser) FirstName() string {
	return u.Claims["firstName"].(string)
}

func (u TestUser) SurName() string {
	return u.Claims["surName"].(string)
}

func (u TestUser) Email() string {
	return u.Claims["email"].(string)
}

func (u TestUser) Username() string {
	return u.Claims["username"].(string)
}

func (u TestUser) IsAdmin() bool {
	return u.Claims["admin"].(bool)
}

func (u TestUser) IsDeveloper() bool {
	return u.Claims["developer"].(bool)
}

func (u TestUser) Set(req *http.Request) error {
	return nil
}

func GetTestUserIdentity() UserIdentity {
	return TestUser{
		Claims: map[string]any{
			"partner": "PARTNER-X",
			"tenant":  []string{"default"},
			"roles": map[string]any{
				"testCase": "customer",
			},
			"firstName": "John",
			"surName":   "Rocket",
			"email":     "j.rocket@example.com",
			"username":  "jrocket",
			"admin":     false,
			"developer": false,
		},
		tenant: "default",
	}
}
