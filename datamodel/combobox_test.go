package datamodel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadTenantCombobox(t *testing.T) {
	path := GetConfigPath("../assets/config", "default", 1)

	tcl, err := loadTenantComboboxList(path, "combobox-prototype.json", 1)
	require.NoError(t, err)

	require.Exactly(t, 2, len(*tcl))
	require.Greater(t, len((*tcl)["inquiry"]), 0)

	merchandiser, ok := (*tcl)["inquiry"]["merchandiser"]
	require.Equal(t, true, ok)
	require.NotNil(t, merchandiser)
	require.Equal(t, "config", *merchandiser.Source)
	require.Equal(t, "merchandiser", *merchandiser.NameInSource)

	team, ok := (*tcl)["inquiry"]["team"]
	require.Equal(t, true, ok)
	require.NotNil(t, team)
	require.Equal(t, false, *team.Translate)
	require.Equal(t, 3, len(team.Content))
}
