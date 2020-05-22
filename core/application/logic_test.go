package application

import (
	"testing"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationUsed(t *testing.T) {
	result, err := GetApplicationUsed("WyJIaUxvIFNpbXVsYXRvciIsIkVuY29kZXIgRGVjb2RlciIsIlFSIENvZGUiLCJETlMgTG9va3VwIiwiUmVxdWVzdCBCaW4iLCJLZWxseSBDcml0ZXJpb24iXQ%3D%3D")

	assert.Empty(t, err)
	assert.NotEmpty(t, result)

	result, err = GetApplicationUsed("")
	assert.Empty(t, err)
	assert.Empty(t, result)
}

func TestGetApplicationUsedFailed(t *testing.T) {
	//url fail
	result, err := GetApplicationUsed("%")
	assert.Empty(t, result)
	assert.NotEmpty(t, err.Error())

	//base64 fail
	result, err = GetApplicationUsed("a")
	assert.Empty(t, result)
	assert.NotEmpty(t, err.Error())

	//json fail
	result, err = GetApplicationUsed("e2ExMjN9")
	assert.Empty(t, result)
	assert.NotEmpty(t, err.Error())
}

func TestGetApplicationsByName(t *testing.T) {
	ret := GetApplicationsByName("dns-lookup")
	assert.NotEmpty(t, ret)

	ret = GetApplicationsByName("asdfadsf")
	assert.Nil(t, ret)
}

func TestGetApplicationWithLiked(t *testing.T) {
	user := &db.User{}
	user.LikedApps = []string{"DNS Lookup"}
	assert.NotEmpty(t, GetApplicationWithLiked(user))
	assert.Empty(t, GetApplicationWithLiked(nil))
}

func TestReloadAppList(t *testing.T) {
	ReloadAppList()
}

func TestAddApplicationUsage(t *testing.T) {
	app := GetAppList()[0].AppCards[0]
	i := app.AmountUsed
	AddApplicationUsage(app)
	j := app.AmountUsed
	assert.Equal(t, j, i+1)
}
