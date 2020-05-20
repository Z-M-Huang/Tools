package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Remove("./test.file.txt")
	os.Exit(ret)
}

func TestRandomString(t *testing.T) {
	length := 50

	s := RandomString(length)

	assert.Equal(t, len(s), length)
}

func TestWriteContentToFile(t *testing.T) {
	content := []string{
		"test",
		"content",
	}
	path := "./test.file.txt"

	err := WriteContentToFile(content, path)

	assert.Empty(t, err)
	_, err = os.Stat(path)
	assert.False(t, os.IsNotExist(err))
}
