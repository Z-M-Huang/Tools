package pages

import (
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//HomePage home page /
func HomePage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	claim := c.Keys[utils.ClaimCtxKey].(*data.JWTClaim)
	if !(claim == nil) {
		user := &dbentity.User{
			Email: claim.Id,
		}
		err := userlogic.Find(utils.DB, user)
		if err == nil {
			if len(user.LikedApps) > 0 {
				response.Data = applicationlogic.GetApplicationWithLiked(user)
			} else {
				response.Data = utils.AppList
			}
		} else {
			response.Data = utils.AppList
		}
	} else {
		response.Data = utils.AppList
	}

	response.Header.Title = "Fun Apps"
	response.Header.Description = "Fun apps, fun personal small projects, and just for fun."
	c.HTML(200, "homepage.gohtml", response)
}
