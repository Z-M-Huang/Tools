package core

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//Home logic
type Home struct{}

//HomePage home page /
func (h *Home) HomePage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	claim := c.Keys[utils.ClaimCtxKey].(*data.JWTClaim)
	if !(claim == nil) {
		user := &db.User{
			Email: claim.Id,
		}
		err := user.Find()
		if err == nil && len(user.LikedApps) > 0 {
			response.Data = webdata.GetApplicationWithLiked(user)
		}
	}

	if response.Data == nil {
		response.Data = webdata.GetAppList()
	}

	response.Header.Title = "Fun Apps"
	response.Header.Description = "Fun apps, fun personal small projects, and just for fun."
	c.HTML(200, "homepage.gohtml", response)
}
