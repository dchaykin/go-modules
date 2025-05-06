package auth

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtTest = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZGV2LWlucXVpcnkuc2V0bG9nLmNvbSJdLCJlTWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJleHAiOjE5MDAxNjY1MzIsImZpcnN0TmFtZSI6IkpvaG5ueSIsImlhdCI6MTc0MjQwMDEzMiwiaXNzIjoiaHR0cHM6Ly9hdXRoLmRldi1pbnF1aXJ5LnNldGxvZy5jb20vYXV0aC9yZWFsbXMvbWFzdGVyIiwicGFydG5lciI6IklOUVVJUlktU09MVVRJT04tTFREIiwicm9sZUNvbmZpZyI6IkN1c3RvbWVyIiwicm9sZUlucXVpcnkiOiJDdXN0b21lciIsInN1YiI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJzdXJOYW1lIjoiUm9ja2V0IiwidXNlck5hbWUiOiJqcm9ja2V0In0.vXMmHGkP4UdfBsBYiF_d8FEGCk_FJesU5c1YU7J5WLI`
const secretTest = `dFwUdN4pCr9kqWNgjCGCYVuL8StRy3sf`
const jwtInvalid = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZGV2LWlucXVpcnkuc2V0bG9nLmNvbSJdLCJlTWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJleHAiOjE5MDAxNjY1MzIsImZpcnN0TmFtZSI6IkpvaG5ueSIsImlhdCI6MTc0MjQwMDEzMiwiaXNzIjoiaHR0cHM6Ly9hdXRoLmRldi1pbnF1aXJ5LnNldGxvZy5jb20vYXV0aC9yZWFsbXMvbWFzdGVyIiwicGFydG5lciI6IklOUVVJUlktU09MVVRJT04tTFREIiwicm9sZUlucXVpcnkiOiJDdXN0b21lciIsInN1YiI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJzdXJOYW1lIjoiUm9ja2V0IiwidXNlck5hbWUiOiJqcm9ja2V0In0.tX5B8_OPOAjN1CHJYF4mKsdZtZj_bh7zmEJHBwdi8xY`

var claimsTest = jwt.MapClaims{
	"iss":         "https://auth.dev-inquiry.setlog.com/auth/realms/master",
	"iat":         1742400132,
	"exp":         1900166532,
	"aud":         []string{"dev-inquiry.setlog.com"},
	"sub":         "jrocket@example.com",
	"firstName":   "Johnny",
	"surName":     "Rocket",
	"eMail":       "jrocket@example.com",
	"roleInquiry": "Customer",
	"roleConfig":  "Customer",
	"partner":     "INQUIRY-SOLUTION-LTD",
	"userName":    "jrocket",
}

func TestValidToken(t *testing.T) {
	user, err := GetUserIdentity(jwtTest, secretTest)
	require.NoError(t, err)
	assert.Equal(t, "INQUIRY-SOLUTION-LTD", user.Partner())
}

func TestInvalidToken(t *testing.T) {
	_, err := GetUserIdentity(jwtInvalid, secretTest)
	assert.EqualError(t, jwt.ErrSignatureInvalid, err.Error())
}

func TestCreateAuthorizationToken(t *testing.T) {
	token, err := CreateAuthorizationToken(claimsTest, secretTest)
	require.NoError(t, err)
	fmt.Println(token)
	require.Equal(t, jwtTest, token)
}
