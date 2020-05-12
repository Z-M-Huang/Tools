package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/logic"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var minPasswordLength int = 12
var googleOauthConfig *oauth2.Config

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("https://%s/google_oauth", utils.Config.Host),
		ClientID:     utils.Config.GoogleOauthConfig.ClientID,
		ClientSecret: utils.Config.GoogleOauthConfig.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

//Login request
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	request := &apidata.LoginRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		WriteResponse(w, response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		WriteResponse(w, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	existingUser := &dbentity.User{
		Email: request.Email,
	}
	err = userlogic.Find(utils.DB, existingUser)
	if err == gorm.ErrRecordNotFound {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "We couldn't find any account for this email address... Maybe you need to create one.",
		})
		WriteResponse(w, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}

	if !utils.ComparePasswords(existingUser.Password, []byte(request.Password)) {
		utils.Logger.Sugar().Errorf("Invalid login attempt received for: %s", request.Email)
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   `Incorrect password. Do you forget your password? If you forget your password, please <a href="#">Click here</a> to reset your password. Uh... We don't have that feature yet, sorry...`,
		})
		WriteResponse(w, response)
		return
	}

	tokenStr, expiresAt, err := userlogic.GenerateJWTToken("Direct Login", request.Email, existingUser.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
	}

	response.Data = true
	logic.SetCookie(w, utils.SessionTokenKey, tokenStr, expiresAt)
	WriteResponse(w, response)
}

//APILogin api login
func APILogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" || !strings.Contains(authHeader, "Basic ") {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}
	authHeader = strings.Replace(authHeader, "Basic ", "", 1)
	decodedAuthBytes, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}

	decodedAuth := string(decodedAuthBytes)
	if !strings.Contains(decodedAuth, ":") {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}

	authSplit := strings.Split(decodedAuth, ":")
	if len(authSplit) != 2 && emailRe.Match([]byte(authSplit[0])) && len(authSplit[1]) < minPasswordLength {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}

	existingUser := &dbentity.User{
		Email: authSplit[0],
	}
	err = userlogic.Find(utils.DB, existingUser)
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("APILogin: Email %s not found", authSplit[0])
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if !utils.ComparePasswords(existingUser.Password, []byte(authSplit[1])) {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}

	tokenStr, expiresAt, err := userlogic.GenerateJWTToken("APILogin", authSplit[0], existingUser.Username, getGravatarLink(authSplit[0], 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	response := &apidata.APILoginResponse{
		TokenType:   "bearer",
		AccessToken: tokenStr,
		ExpiresIn:   int64(expiresAt.Sub(time.Now().UTC()).Seconds()),
	}

	jsonBody, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}

//SignUp request
func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	request := &apidata.CreateAccountRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		WriteResponse(w, response)
		return
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid login request.",
		})
		WriteResponse(w, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if !emailRe.Match([]byte(request.Email)) {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid email address.",
		})
		WriteResponse(w, response)
		return
	}

	if request.ConfirmPassword != request.Password {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Password doesn't match",
		})
		WriteResponse(w, response)
		return
	}

	if len(request.Password) < minPasswordLength {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   fmt.Sprintf("Password has minimum length of %d characters.", minPasswordLength),
		})
		WriteResponse(w, response)
		return
	}

	existingUser := &dbentity.User{
		Email: request.Email,
	}
	err = userlogic.Find(utils.DB, existingUser)
	if err == nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Email address already exists, please try to remember the password, since password recovery function is not yet built. If you cant remember your password, good luck... The password is hashed, and even as an admin, I have no clue what's your password could be... See ya.",
		})
		WriteResponse(w, response)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}

	existingUser = &dbentity.User{
		Username: request.Username,
	}
	err = userlogic.Find(utils.DB, existingUser)
	if err == nil {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Username already taken. Can't you think of something else? Try harder",
		})
		WriteResponse(w, response)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}

	user := &dbentity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: utils.HashAndSalt([]byte(request.Password)),
	}
	err = userlogic.Save(utils.DB, user)
	if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}

	tokenStr, expiresAt, err := userlogic.GenerateJWTToken("Direct Login", request.Email, user.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
	}

	logic.SetCookie(w, utils.SessionTokenKey, tokenStr, expiresAt)
	response.Data = true
	WriteResponse(w, response)
	return
}

//UpdatePassword api
func UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	request := &apidata.UpdatePasswordRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid sign up request.",
		})
		WriteResponse(w, response)
		return
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid sign up request.",
		})
		WriteResponse(w, response)
		return
	}

	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)

	if request.Password != request.ConfirmPassword {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Password doesn't match.",
		})
		WriteResponse(w, response)
		return
	} else if len(request.Password) < minPasswordLength {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   fmt.Sprintf("Password has minimum length of %d.", minPasswordLength),
		})
		WriteResponse(w, response)
		return
	}

	dbUser := &dbentity.User{
		Email: claim.Id,
	}
	err = userlogic.Find(utils.DB, dbUser)
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("User not found for %s in UpdatePassword", claim.Id)
		WriteUnexpectedError(w, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}

	if dbUser.Password != "" && !utils.ComparePasswords(dbUser.Password, []byte(request.CurrentPassword)) {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Current password is different compared to what's in database... Try harder...",
		})
		WriteResponse(w, response)
		return
	} else if dbUser.Password != "" && utils.ComparePasswords(dbUser.Password, []byte(request.Password)) {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "New password is exactly the same as the old password...",
		})
		WriteResponse(w, response)
		return
	}

	dbUser.Password = utils.HashAndSalt([]byte(request.Password))
	err = userlogic.Save(utils.DB, dbUser)

	if err != nil {
		utils.Logger.Error(err.Error())
		WriteUnexpectedError(w, response)
		return
	}
	response.SetAlert(&data.AlertData{
		IsSuccess: true,
		Message:   "Password is updated.",
	})
	WriteResponse(w, response)
}

//GoogleLogin google login request
func GoogleLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	state := utils.RandomString(10)
	err := utils.RedisClient.Set(state, "1", 20*time.Minute).Err()
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleCallback handle google callback
func GoogleCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getGoogleUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	dbUser := &dbentity.User{
		Email: user.Email,
	}
	err = userlogic.Find(utils.DB, dbUser)
	if err == gorm.ErrRecordNotFound {
		//User not found
		dbUser.Username = user.Name
		dbUser.GoogleID = user.ID
		dbUser.Email = user.Email
		err = userlogic.Save(utils.DB, dbUser)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	} else if err != nil {
		//Something else happend
		utils.Logger.Error(err.Error())
	}

	tokenStr, expiresAt, err := userlogic.GenerateJWTToken("Google", user.Email, user.Name, user.Picture)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		logic.SetCookie(w, utils.SessionTokenKey, tokenStr, expiresAt)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getGoogleUserInfo(state string, code string) (*apidata.GoogleUserInfo, error) {
	_, err := utils.RedisClient.Get(state).Result()
	if err != nil {
		return nil, fmt.Errorf("invalid oauth state")
	}
	utils.RedisClient.Del(state)
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
	user := &apidata.GoogleUserInfo{}
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
