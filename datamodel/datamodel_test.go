package datamodel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTenantConfig(t *testing.T) {
	tc, err := LoadDataModelByRole("testdata", "customer", 1)
	require.NoError(t, err)

	require.Nil(t, tc.Roles)
	require.NotNil(t, tc.Cmbs)

	require.Exactly(t, 2, len(*tc.Cmbs))
	require.Greater(t, len((*tc.Cmbs)["roles"]), 0)

	userPartner, ok := (*tc.Cmbs)["user"]["partner"]
	require.Equal(t, true, ok)
	require.NotNil(t, userPartner)
	require.Equal(t, "/app-config/api/partner/cmbs/userPartner", *userPartner.Source)
	require.Equal(t, "userPartner", userPartner.Name)

	roleNames, ok := (*tc.Cmbs)["roles"]["name"]
	require.Equal(t, true, ok)
	require.NotNil(t, roleNames)
	require.Equal(t, false, *roleNames.Translate)
	require.Equal(t, 3, len(roleNames.Content))
}
