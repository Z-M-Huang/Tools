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

func TestSearchAppListByNames(t *testing.T) {
	ret := SearchAppListByNames([]string{"dns-lookup", "kelly-criterion"})
	assert.NotEmpty(t, ret)

	ret = SearchAppListByNames([]string{"asdfasdfasdf"})
	assert.Empty(t, ret)
}

func TestSearchAppListByNamesWithLikes(t *testing.T) {
	likedApps := []string{"dns-lookup", "kelly-criterion"}
	user := &db.User{}
	user.LikedApps = []string{"DNS Lookup", "Kelly Criterion"}
	ret := SearchAppListByNamesWithLikes(user, likedApps)
	assert.NotEmpty(t, ret)

	ret = SearchAppListByNamesWithLikes(user, []string{"asdfasdfasdf"})
	assert.Empty(t, ret)
}

func TestGetApplicationWithLiked(t *testing.T) {
	user := &db.User{}
	user.LikedApps = []string{"DNS Lookup"}
	assert.NotEmpty(t, GetAppListWithLiked(user))
	assert.Empty(t, GetAppListWithLiked(nil))
}

func TestReloadAppList(t *testing.T) {
	ReloadAppList()
}

func TestLoadSearchMappings(t *testing.T) {
	LoadSearchMappings()
}

func TestAddApplicationUsage(t *testing.T) {
	app := GetAppList()[0].AppCards[0]
	i := app.AmountUsed
	AddApplicationUsage(app)
	j := app.AmountUsed
	assert.Equal(t, j, i+1)
}
