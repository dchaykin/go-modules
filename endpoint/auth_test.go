package endpoint

import (
	"os"
	"testing"

	"github.com/dchaykin/go-modules/helper"
	"github.com/stretchr/testify/require"
)

func TestCreateTokenByCredentials(t *testing.T) {
	helper.LoadAccessData("../.do-not-commit/env.vars")
	userIdentity, err := CreateUserIdentityByCredentials(os.Getenv("TECH_USER"), os.Getenv("TECH_PASS"), os.Getenv("AUTH_SECRET"))
	require.NoError(t, err)
	require.NotNil(t, userIdentity)
	require.Equal(t, "Cycle", userIdentity.FirstName())
	require.Equal(t, "Bot", userIdentity.SurName())
}
