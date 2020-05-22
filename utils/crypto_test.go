package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndSalt(t *testing.T) {
	passToHash := "asidfjopiawje2341234!!!"

	hashed := HashAndSalt([]byte(passToHash))

	assert.NotEqual(t, passToHash, hashed)
	assert.NotEmpty(t, hashed)
	assert.True(t, ComparePasswords(hashed, []byte(passToHash)))
}
