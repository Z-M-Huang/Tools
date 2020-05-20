package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAlert(t *testing.T) {
	errorMessage := "test"
	p := &PageResponse{}

	p.SetAlert(&AlertData{
		Message: errorMessage,
	})

	assert.NotEmpty(t, p.Header.Alert)
	assert.Equal(t, p.Header.Alert.Message, errorMessage)
}

func TestSetLogin(t *testing.T) {
	username := "test"
	p := &PageResponse{}

	p.SetLogin(&LoginData{
		Username: username,
	})

	assert.NotEmpty(t, p.Header.Nav.Login)
	assert.Equal(t, p.Header.Nav.Login.Username, username)
}

func TestSetNavStyleName(t *testing.T) {
	style := "default"
	link := "https://localhost"
	integrity := "anyrandomnumber"
	p := &PageResponse{}

	p.SetNavStyleName(&PageStyleData{
		Name:      style,
		Link:      link,
		Integrity: integrity,
	})

	assert.NotEmpty(t, p.Header.Nav.StyleName)
	assert.NotEmpty(t, p.Header.PageStyle)
	assert.Equal(t, p.Header.PageStyle.Name, style)
	assert.Equal(t, p.Header.PageStyle.Link, link)
	assert.Equal(t, p.Header.PageStyle.Integrity, integrity)
}
