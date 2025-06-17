package auth

import (
	"fmt"
	"net/http"
)

type TestUser struct {
	Claims        map[string]any
	CurrentTenant string
}

func (u TestUser) Partner() string {
	return u.Claims["partner"].(string)
}

func (u TestUser) Tenant() string {
	return u.CurrentTenant
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
	return u.Claims["eMail"].(string)
}

func (u TestUser) Username() string {
	return u.Claims["userName"].(string)
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

func GetTestUserIdentity() TestUser {
	return TestUser{
		Claims: map[string]any{
			"partner": "PARTNER-X",
			"tenant":  []string{"default"},
			"roles": map[string]any{
				"testCase": "customer",
			},
			"firstName": "John",
			"surName":   "Rocket",
			"eMail":     "j.rocket@example.com",
			"userName":  "jrocket",
			"admin":     false,
			"developer": false,
		},
		CurrentTenant: "default",
	}
}
