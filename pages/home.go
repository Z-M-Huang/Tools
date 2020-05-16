package pages

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/constval"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/gin-gonic/gin"
)

//HomePage home page /
func HomePage(c *gin.Context) {
	response := c.Keys[constval.ResponseCtxKey].(*data.Response)
	claim := c.Keys[constval.ClaimCtxKey].(*data.JWTClaim)
	if !(claim == nil) {
		user := &db.User{
			Email: claim.Id,
		}
		err := user.Find()
		if err == nil {
			if len(user.LikedApps) > 0 {
				response.Data = webdata.GetApplicationWithLiked(user)
			} else {
				response.Data = webdata.AppList
			}
		} else {
			response.Data = webdata.AppList
		}
	} else {
		response.Data = webdata.AppList
	}

	response.Header.Title = "Fun Apps"
	response.Header.Description = "Fun apps, fun personal small projects, and just for fun."
	c.HTML(200, "homepage.gohtml", response)
}
