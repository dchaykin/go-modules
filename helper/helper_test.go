package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatFromString(t *testing.T) {
	require.EqualValues(t, 0, FloatFromString(""))
	require.EqualValues(t, 6.4, FloatFromString("6.4"))
	require.EqualValues(t, 1.0, FloatFromString("1"))
	require.EqualValues(t, 0, FloatFromString("6,4"))
}
