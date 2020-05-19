package home

import (
	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/core/application"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/gin-gonic/gin"
)

//Page home page
type Page struct{}

//Home home page /
func (Page) Home(c *gin.Context) {
	response := core.GetResponseInContext(c.Keys)
	claim := account.GetClaimInContext(c.Keys)
	if !(claim == nil) {
		user := &db.User{
			Email: claim.Id,
		}
		err := user.Find()
		if err == nil && len(user.LikedApps) > 0 {
			response.Data = application.GetApplicationWithLiked(user)
		}
	}

	if response.Data == nil {
		response.Data = application.GetAppList()
	}

	response.Header.Title = "Fun Apps"
	response.Header.Description = "Fun apps, fun personal small projects, and just for fun."
	c.HTML(200, "homepage.gohtml", response)
}
