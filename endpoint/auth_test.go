package endpoint

import (
	"os"
	"testing"

	"github.com/dchaykin/go-modules/helper"
	"github.com/stretchr/testify/require"
)

func TestCreateTokenByCredentials(t *testing.T) {
	helper.LoadAccessData("../.do-not-commit/env.vars")
	token, err := CreateTokenByCredentials(os.Getenv("TECH_USER"), os.Getenv("TECH_PASS"))
	require.NoError(t, err)
	require.NotNil(t, token)
}
