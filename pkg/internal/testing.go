package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ErrorsMatch(t *testing.T, expected, actual error) {
	if expected != nil {
		assert.EqualError(t, actual, expected.Error())
	} else {
		assert.Nil(t, actual)
	}
}
