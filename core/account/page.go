package account

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//Page account page
type Page struct{}

//Signup /signup
func (Page) Signup(c *gin.Context) {
	claim := core.GetClaimInContext(c.Keys)
	if claim == nil {
		response := core.GetResponseInContext(c.Keys)
		response.Header.Title = "Signup - Fun Apps"
		response.Header.Description = "Signup - create an account"
		c.HTML(200, "signup.gohtml", response)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

//Login /login
func (Page) Login(c *gin.Context) {
	claim := core.GetClaimInContext(c.Keys)
	if claim != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		response := core.GetResponseInContext(c.Keys)
		response.Header.Title = "Login - Fun Apps"
		response.Header.Description = "Login"
		redirectURL, ok := c.Request.URL.Query()["redirect"]
		if ok && len(redirectURL) > 0 {
			response.SetAlert(&data.AlertData{
				IsDanger: true,
				Message:  "Please login first.",
			})
		}
		c.HTML(200, "login.gohtml", response)
	}
}

//Account /account requires claim
func (Page) Account(c *gin.Context) {
	claim := core.GetClaimInContext(c.Keys)
	if claim != nil {
		response := core.GetResponseInContext(c.Keys)
		response.Header.Title = "Account - Fun Apps"
		response.Header.Description = "Manage account"
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
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login?redirect=/account")
	}
}
