package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testValidatePayload struct {
	Name string `validate:"required"`
}

func TestValidate(t *testing.T) {
	require.NoError(t, Validate(testValidatePayload{Name: "ok"}))
	require.Error(t, Validate(testValidatePayload{}))
}
