package application

import (
	"testing"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/stretchr/testify/assert"
)

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
