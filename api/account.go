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
	response := &data.Response{}
	request := &apidata.LoginRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid login request."
		writeResponse(w, response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid login request."
		writeResponse(w, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	existingUser := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: request.Email,
	}).First(&existingUser); db.RecordNotFound() || existingUser == nil {
		response.Alert.IsDanger = true
		response.Alert.Message = "We couldn't find any account for this email address... Maybe you need to create one"
		writeResponse(w, response)
		return
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		writeUnexpectedError(w, response)
		return
	}

	if !utils.ComparePasswords(existingUser.Password, []byte(request.Password)) {
		utils.Logger.Sugar().Errorf("Invalid login attempt received for: %s", request.Email)
		response.Alert.IsWarning = true
		response.Alert.Message = `Incorrect password. Do you forget your password? If you forget your password, please <a href="#">Click here</a> to reset your password. Uh... We don't have that feature yet, sorry...`
		writeResponse(w, response)
		return
	}

	tokenStr, expiresAt, err := utils.GenerateJWTToken("Direct Login", request.Email, existingUser.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		writeUnexpectedError(w, response)
	}

	SetAuthCookie(w, tokenStr, expiresAt)
	writeResponse(w, response)
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

	existingUser := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: authSplit[0],
	}).First(&existingUser); db.RecordNotFound() || existingUser == nil {
		utils.Logger.Sugar().Errorf("APILogin: Email %s not found", authSplit[0])
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if !utils.ComparePasswords(existingUser.Password, []byte(authSplit[1])) {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return
	}

	tokenStr, expiresAt, err := utils.GenerateJWTToken("APILogin", authSplit[0], existingUser.Username, getGravatarLink(authSplit[0], 50))
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
	response := &data.Response{}
	request := &apidata.CreateAccountRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid sign up request."
		writeResponse(w, response)
		return
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid sign up request."
		writeResponse(w, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if !emailRe.Match([]byte(request.Email)) {
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid email address."
		writeResponse(w, response)
		return
	}

	if request.ConfirmPassword != request.Password {
		response.Alert.IsWarning = true
		response.Alert.Message = "Password doesn't match"
		writeResponse(w, response)
		return
	}

	if len(request.Password) < minPasswordLength {
		response.Alert.IsWarning = true
		response.Alert.Message = fmt.Sprintf("Password has minimum length of %d characters.", minPasswordLength)
		writeResponse(w, response)
		return
	}

	existingUser := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: request.Email,
	}).First(&existingUser); !db.RecordNotFound() {
		response.Alert.IsWarning = true
		response.Alert.Message = "Email address already exists, please try to remember the password, since password recovery function is not yet built. If you cant remember your password, good luck... The password is hashed, and even as an admin, I have no clue what's your password could be... See ya."
		writeResponse(w, response)
		return
	}

	if db := utils.DB.Where(dbentity.User{
		Username: request.Username,
	}).First(&existingUser); !db.RecordNotFound() {
		response.Alert.IsWarning = true
		response.Alert.Message = "Username already taken. Can't you think of something else? Try harder"
		writeResponse(w, response)
		return
	} else if db.RecordNotFound() {

	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		writeUnexpectedError(w, response)
		return
	}

	user := &dbentity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: utils.HashAndSalt([]byte(request.Password)),
	}
	if db := utils.DB.Save(user).Scan(&user); db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		writeUnexpectedError(w, response)
		return
	}

	tokenStr, expiresAt, err := utils.GenerateJWTToken("Direct Login", request.Email, user.Username, getGravatarLink(request.Email, 50))
	if err != nil {
		utils.Logger.Error(err.Error())
		writeUnexpectedError(w, response)
	}

	SetAuthCookie(w, tokenStr, expiresAt)
	response.Data = true
	writeResponse(w, response)
	return
}

//UpdatePassword api
func UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := &data.Response{}
	request := &apidata.UpdatePasswordRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid sign up request."
		writeResponse(w, response)
		return
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "Invalid sign up request."
		writeResponse(w, response)
		return
	}

	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)

	if request.Password != request.ConfirmPassword {
		response.Alert.IsWarning = true
		response.Alert.Message = "Password doesn't match."
		writeResponse(w, response)
		return
	} else if len(request.Password) < minPasswordLength {
		response.Alert.IsWarning = true
		response.Alert.Message = fmt.Sprintf("Password has minimum length of %d.", minPasswordLength)
		writeResponse(w, response)
		return
	}

	dbUser := &dbentity.User{
		Email: claim.Id,
	}
	if db := utils.DB.Where(dbUser).First(&dbUser); db.RecordNotFound() {
		utils.Logger.Sugar().Errorf("User not found for %s", claim.Id)
		writeUnexpectedError(w, response)
		return
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		writeUnexpectedError(w, response)
		return
	}

	if !utils.ComparePasswords(dbUser.Password, []byte(request.CurrentPassword)) {
		response.Alert.IsWarning = true
		response.Alert.Message = "Current password is different compared to what's in database... Try harder..."
		writeResponse(w, response)
		return
	}

	if utils.ComparePasswords(dbUser.Password, []byte(request.Password)) {
		response.Alert.IsWarning = true
		response.Alert.Message = "New password is exactly the same as the old password..."
		writeResponse(w, response)
		return
	}

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		dbUser.Password = utils.HashAndSalt([]byte(request.Password))

		if db := tx.Save(dbUser).Scan(&dbUser); db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to update user with new password %s", db.Error.Error()))
		}
		return nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		writeUnexpectedError(w, response)
		return
	}
	response.Alert.IsSuccess = true
	response.Alert.Message = "Password is updated."
	writeResponse(w, response)
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
	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		dbUser := &dbentity.User{
			Email: user.Email,
		}
		if db := tx.Where(dbUser).First(&dbUser); db.RecordNotFound() {
			dbUser.Username = user.Name
			dbUser.GoogleID = user.ID
			dbUser.Email = user.Email
			if db = tx.Save(dbUser).Scan(&dbUser); db.Error != nil {
				return fmt.Errorf(fmt.Sprintf("failed to save new user in GoogleCallBack %s", db.Error))
			}
		} else if db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to save new user in GoogleCallBack %s", db.Error))
		}
		return nil
	})
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	tokenStr, expiresAt, err := utils.GenerateJWTToken("Google", user.Email, user.Name, user.Picture)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		SetAuthCookie(w, tokenStr, expiresAt)
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

//SetAuthCookie set auth cookie
func SetAuthCookie(w http.ResponseWriter, tokenStr string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:       utils.SessionTokenKey,
		Value:      tokenStr,
		Path:       "/",
		Domain:     utils.Config.Host,
		Expires:    expiresAt,
		RawExpires: expiresAt.String(),
	})
}

func getGravatarLink(email string, size uint) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}
