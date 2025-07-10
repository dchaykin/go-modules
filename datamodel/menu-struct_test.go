package datamodel

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMenuConfig_ReadFromFile(t *testing.T) {
	mc := MenuConfig{}
	err := mc.ReadFromFile("testdata-001")
	require.NoError(t, err)

	menu := mc.CreateMenuByRole("customer")

	jsonData, _ := json.MarshalIndent(menu, "", "  ")
	fmt.Println(string(jsonData))

	require.EqualValues(t, len(menu), 3)
	require.EqualValues(t, "dashboard", menu[0].Name)
	require.EqualValues(t, "partner", menu[1].Name)
	require.EqualValues(t, "settings", menu[2].Name)

	require.EqualValues(t, 0, len(menu[0].SubItems))
	require.EqualValues(t, 3, len(menu[1].SubItems))
	require.EqualValues(t, "overview", menu[1].SubItems[0].Name)
	require.EqualValues(t, "input", menu[1].SubItems[1].Name)
	require.EqualValues(t, "search", menu[1].SubItems[2].Name)

	require.EqualValues(t, 1, len(menu[2].SubItems))
	require.EqualValues(t, "preferences", menu[2].SubItems[0].Name)
}
