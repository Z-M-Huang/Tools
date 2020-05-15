package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//SignupPage /signup
func SignupPage(c *gin.Context) {
	if c.Keys[utils.ClaimCtxKey].(*data.JWTClaim) == nil {
		response := c.Keys[utils.ResponseCtxKey].(*data.Response)
		response.Header.Title = "Signup - Fun Apps"
		response.Header.Description = "Signup - create an account"
		utils.Templates.ExecuteTemplate(c.Writer, "signup.gohtml", response)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

//LoginPage /login
func LoginPage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	response.Header.Title = "Login - Fun Apps"
	response.Header.Description = "Login"
	redirectURL, ok := c.Request.URL.Query()["redirect"]
	if ok && len(redirectURL) > 0 {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Please login first.",
		})
	}
	if c.Keys[utils.ClaimCtxKey].(*data.JWTClaim) == nil {
		utils.Templates.ExecuteTemplate(c.Writer, "login.gohtml", response)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

//AccountPage /account requires claim
func AccountPage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	response.Header.Title = "Account - Fun Apps"
	response.Header.Description = "Manage account"
	claim := c.Keys[utils.ClaimCtxKey].(*data.JWTClaim)
	user := &dbentity.User{
		Email: claim.Id,
	}
	err := userlogic.Find(utils.DB, user)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
		})
	} else {
		response.Data = webdata.AccountPageData{
			HasPassword: user.Password != "",
		}
	}
	utils.Templates.ExecuteTemplate(c.Writer, "account.gohtml", response)
}
