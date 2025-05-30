package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapError(t *testing.T) {
	err := errorFunc3()
	require.Error(t, err)

	Error(err.Error())
}

func errorFunc3() error {
	return WrapError(errorFunc2())
}

func errorFunc2() error {
	return WrapError(errorFunc1())
}

func errorFunc1() error {
	return WrapError(fmt.Errorf("error occured"))
}
