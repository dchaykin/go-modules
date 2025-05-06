package datamodel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndExtractUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	require.NoError(t, err)
	require.Equal(t, 32, len(uuid))

	extractedTime, err := ExtractTimeFromUUID(uuid)
	require.NoError(t, err)

	now := time.Now()
	diff := now.Sub(extractedTime)
	require.LessOrEqual(t, diff, 2*time.Second)
	require.GreaterOrEqual(t, diff, -2*time.Second)
}
