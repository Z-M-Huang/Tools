package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/constval"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//SignupPage /signup
func SignupPage(c *gin.Context) {
	if c.Keys[constval.ClaimCtxKey].(*data.JWTClaim) == nil {
		response := c.Keys[constval.ResponseCtxKey].(*data.Response)
		response.Header.Title = "Signup - Fun Apps"
		response.Header.Description = "Signup - create an account"
		c.HTML(200, "signup.gohtml", response)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

//LoginPage /login
func LoginPage(c *gin.Context) {
	response := c.Keys[constval.ResponseCtxKey].(*data.Response)
	response.Header.Title = "Login - Fun Apps"
	response.Header.Description = "Login"
	redirectURL, ok := c.Request.URL.Query()["redirect"]
	if ok && len(redirectURL) > 0 {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Please login first.",
		})
	}
	if c.Keys[constval.ClaimCtxKey].(*data.JWTClaim) == nil {
		c.HTML(200, "login.gohtml", response)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

//AccountPage /account requires claim
func AccountPage(c *gin.Context) {
	response := c.Keys[constval.ResponseCtxKey].(*data.Response)
	response.Header.Title = "Account - Fun Apps"
	response.Header.Description = "Manage account"
	claim := c.Keys[constval.ClaimCtxKey].(*data.JWTClaim)
	user := &db.User{
		Email: claim.Id,
	}
	err := user.Find()
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
	c.HTML(200, "account.gohtml", response)
}
