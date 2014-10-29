package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinOfTwoIntegers(t *testing.T) {
	assert.Equal(t, 3, MinI(3, 4))
	assert.Equal(t, 3, MinI(4, 3))
	assert.Equal(t, 3, MinI(3, 3))
}

func TestMinOfTwoFloats(t *testing.T) {
	assert.Equal(t, 3.1, MinF(3.1, 4.2))
	assert.Equal(t, 3.1, MinF(4.2, 3.1))
	assert.Equal(t, 3.1, MinF(3.1, 3.1))
}

func TestRoundingDown(t *testing.T) {
	assert.Equal(t, 3, Round(3.1))
}

func TestRoundingUp(t *testing.T) {
	assert.Equal(t, 4, Round(3.51))
}
