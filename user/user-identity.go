package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/dchaykin/mygolib/auth"
	"github.com/dchaykin/mygolib/log"
)

type UserIdentity interface {
	auth.SimpleUserIdentity
	Partner() string
	Tenant() string
	RoleByApp(appName string) string
	Apps() []string
}

type userToken struct {
	auth.UserClaims
	CurrentTenant string `json:"currentTenant"`
}

func (j userToken) Partner() string {
	claim, ok := j.Claims["partner"]
	if !ok {
		log.Warn("Claim 'partner' not found")
		return ""
	}
	return claim.(string)
}

func (j userToken) RoleByApp(appName string) string {
	rolesClaim, ok := j.Claims["roles"]
	if !ok {
		log.Warn("User has no roles")
		return ""
	}
	roles := rolesClaim.(map[string]any)
	if result, ok := roles[appName]; ok {
		return fmt.Sprintf("%v", result)
	}
	log.Warn("User has no role for %s. Available roles: %v", appName, rolesClaim)
	return ""
}

func (j userToken) Apps() []string {
	rolesClaim, ok := j.Claims["roles"]
	if !ok {
		log.Warn("User has no roles")
		return nil
	}
	apps := []string{}
	roles := rolesClaim.(map[string]any)
	for app := range roles {
		apps = append(apps, fmt.Sprintf("%v", app))
	}
	return apps
}

func (j userToken) tenantList() []string {
	claim, ok := j.Claims["tenant"]
	if !ok {
		log.Warn("No tenant found")
		return []string{}
	}
	result := []string{}
	for _, v := range claim.([]any) {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return result
}

func (j userToken) Tenant() string {
	if j.CurrentTenant == "" {
		return "default"
	}
	if slices.Contains(j.tenantList(), j.CurrentTenant) {
		return j.CurrentTenant
	}
	return j.CurrentTenant
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
	authorization, err := auth.CreateAuthorizationToken(j.Claims, os.Getenv("AUTH_SECRET"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(authorization))
	return nil
}
