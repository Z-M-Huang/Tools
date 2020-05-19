package account

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//API api endpoints
type API struct{}

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var minPasswordLength int = 12
var googleOauthConfig *oauth2.Config

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("https://%s/google_oauth", data.Config.Host),
		ClientID:     data.Config.GoogleOauthConfig.ClientID,
		ClientSecret: data.Config.GoogleOauthConfig.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	if !data.Config.HTTPS {
		googleOauthConfig.RedirectURL = fmt.Sprintf("http://%s/google_oauth", data.Config.Host)
	}
}

//Login request
func (API) Login(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*core.Response)
	request := &LoginRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	existingUser := &db.User{
		Email: request.Email,
	}
	err = existingUser.Find()
	if err == gorm.ErrRecordNotFound {
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "We couldn't find any account for this email address... Maybe you need to create one.",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}

	if !utils.ComparePasswords(existingUser.Password, []byte(request.Password)) {
		utils.Logger.Sugar().Errorf("Invalid login attempt received for: %s", request.Email)
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   `Incorrect password. Do you forget your password? If you forget your password, please <a href="#">Click here</a> to reset your password. Uh... We don't have that feature yet, sorry...`,
		})
		core.WriteResponse(c, 400, response)
		return
	}

	tokenStr, expiresAt, err := GenerateJWTToken("Direct Login", request.Email, existingUser.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
	}

	uri, err := url.ParseRequestURI(c.GetHeader("Referer"))
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	result := &LoginResponse{
		IsSuccess: true,
	}

	redirect := uri.Query().Get("redirect")
	if len(redirect) > 0 {
		uri, err = url.ParseRequestURI(redirect)
		if err == nil && uri.Host == "" {
			result.Redirect = redirect
		} else {
			utils.Logger.Sugar().Warnf("Illegal redirect uri detected. %s", uri.RequestURI)
		}
	}
	response.Data = result
	core.SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, true)
	core.WriteResponse(c, 200, response)
}

//Logout logout
func (API) Logout(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*core.Response)
	core.SetCookie(c, utils.SessionCookieKey, "", time.Now().AddDate(-10, 1, 1), true)
	response.Data = true
	core.WriteResponse(c, 200, response)
}

//SignUp request
func (API) SignUp(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*core.Response)
	request := &CreateAccountRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if !emailRe.Match([]byte(request.Email)) {
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid email address.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if request.ConfirmPassword != request.Password {
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Password doesn't match",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if len(request.Password) < minPasswordLength {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   fmt.Sprintf("Password has minimum length of %d characters.", minPasswordLength),
		})
		core.WriteResponse(c, 400, response)
		return
	}

	existingUser := &db.User{
		Email: request.Email,
	}
	err = existingUser.Find()
	if err == nil {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Email address already exists, please try to remember the password, since password recovery function is not yet built. If you cant remember your password, good luck... The password is hashed, and even as an admin, I have no clue what's your password could be... See ya.",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}

	existingUser = &db.User{
		Username: request.Username,
	}
	err = existingUser.Find()
	if err == nil {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Username already taken. Can't you think of something else? Try harder",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}

	user := &db.User{
		Username: request.Username,
		Email:    request.Email,
		Password: utils.HashAndSalt([]byte(request.Password)),
	}
	err = user.Save()
	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}

	tokenStr, expiresAt, err := GenerateJWTToken("Direct Login", request.Email, user.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
	}

	core.SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, true)
	response.Data = true
	core.WriteResponse(c, 200, response)
	return
}

//UpdatePassword api
func (API) UpdatePassword(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*core.Response)
	request := &UpdatePasswordRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid sign up request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	claim := c.Keys[utils.ClaimCtxKey].(*JWTClaim)

	if request.Password != request.ConfirmPassword {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Password doesn't match.",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if len(request.Password) < minPasswordLength {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   fmt.Sprintf("Password has minimum length of %d.", minPasswordLength),
		})
		core.WriteResponse(c, 400, response)
		return
	}

	dbUser := &db.User{
		Email: claim.Id,
	}
	err = dbUser.Find()
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("User not found for %s in UpdatePassword", claim.Id)
		core.WriteUnexpectedError(c, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}

	if dbUser.Password != "" && !utils.ComparePasswords(dbUser.Password, []byte(request.CurrentPassword)) {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Current password is different compared to what's in database... Try harder...",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if dbUser.Password != "" && utils.ComparePasswords(dbUser.Password, []byte(request.Password)) {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "New password is exactly the same as the old password...",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	dbUser.Password = utils.HashAndSalt([]byte(request.Password))
	err = dbUser.Save()

	if err != nil {
		utils.Logger.Error(err.Error())
		core.WriteUnexpectedError(c, response)
		return
	}
	response.SetAlert(&core.AlertData{
		IsSuccess: true,
		Message:   "Password is updated.",
	})
	core.WriteResponse(c, 400, response)
}

//GoogleLogin google login request
func (API) GoogleLogin(c *gin.Context) {
	state := utils.RandomString(10)
	err := db.RedisSet(state, "1", 20*time.Minute)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//GoogleCallback handle google callback
func (API) GoogleCallback(c *gin.Context) {
	user, err := getGoogleUserInfo(c.Request.FormValue("state"), c.Request.FormValue("code"))
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	dbUser := &db.User{
		Email: user.Email,
	}
	err = dbUser.Find()
	if err == gorm.ErrRecordNotFound {
		//User not found
		dbUser.Username = user.Name
		dbUser.GoogleID = user.ID
		dbUser.Email = user.Email
		err = dbUser.Save()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	} else if err != nil {
		//Something else happend
		utils.Logger.Error(err.Error())
	}

	tokenStr, expiresAt, err := GenerateJWTToken("Google", user.Email, user.Name, user.Picture)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		core.SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, false)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getGoogleUserInfo(state string, code string) (*GoogleUserInfo, error) {
	if !db.RedisExist(state) {
		return nil, fmt.Errorf("invalid oauth state")
	}
	db.RedisDelete(state)
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	user := &GoogleUserInfo{}
	err = json.Unmarshal(contents, &user)
	if err != nil {
		return nil, fmt.Errorf("failed parsing response to struct %s", err.Error())
	}
	return user, nil
}

func getGravatarLink(email string, size uint) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}
