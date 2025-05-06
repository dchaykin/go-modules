package datamodel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureUUID(t *testing.T) {
	doc := map[string]interface{}{}

	err := EnsureUUID(doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc["uuid"].(string)))

	doc["uuid"] = "invalid-uuid"
	err = EnsureUUID(doc)
	require.NoError(t, err)
	require.EqualValues(t, 32, len(doc["uuid"].(string)))

	uuid, err := GenerateUUID()
	doc["uuid"] = uuid
	require.NoError(t, err)
	err = EnsureUUID(doc)
	require.NoError(t, err)
	require.EqualValues(t, uuid, doc["uuid"].(string))
}
