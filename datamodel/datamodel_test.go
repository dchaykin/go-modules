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

func TestCustomFieldValueByType(t *testing.T) {
	field := CustomField{
		"type": FieldTypeDate,
	}

	require.Equal(t, "2025-06-16", field.ValueByType("2025-06-16T00:00:00.000+02:00"))

	field["type"] = FieldTypeDateTime
	require.Equal(t, "2025-06-16T00:00:00.000+02:00", field.ValueByType("2025-06-16T00:00:00.000+02:00"))

	field["type"] = FieldTypeDate
	require.Equal(t, "2025", field.ValueByType("2025"))
}
