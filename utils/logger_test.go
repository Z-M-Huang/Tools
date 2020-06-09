package utils

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	err := errors.New("Test")
	Logger.Error(err.Error())
}
