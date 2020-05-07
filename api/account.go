package api

import (
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

	//TODO implement some avatar service...
	tokenStr, expiresAt, err := utils.GenerateJWTToken("Direct Login", request.Email, existingUser.Username, "")
	if err != nil {
		utils.Logger.Error(err.Error())
		writeUnexpectedError(w, response)
	}

	SetAuthCookie(w, tokenStr, expiresAt)
	writeResponse(w, response)
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

	if len(request.Password) < 12 {
		response.Alert.IsWarning = true
		response.Alert.Message = "Password has minimum length of 12 characters."
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

	//TODO implement some avatar service...
	tokenStr, expiresAt, err := utils.GenerateJWTToken("Direct Login", request.Email, user.Username, "")
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
	resp := &data.Response{}

	token := r.Header.Get("Authorization")
	if token == "" {
		resp.Alert.IsDanger = true
		resp.Alert.Message = "Unauthorized"
		resp.Data = false
		writeResponse(w, resp)
		return
	}

	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse request data", http.StatusBadRequest)
		return
	}

	password := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		http.Error(w, "Password doesn't match.", http.StatusBadRequest)
		return
	} else if len(password) < 12 {
		http.Error(w, "Password has minimum length of 12", http.StatusBadRequest)
		return
	}

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		dbUser := &dbentity.User{
			Email: claim.Id,
		}
		if db := tx.Where(dbUser).First(&dbUser); db.RecordNotFound() {
			return fmt.Errorf("user not found for email: %s", claim.Id)
		} else if db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to user in UpdatePassword %s", db.Error.Error()))
		}

		dbUser.Password = utils.HashAndSalt([]byte(password))

		if db := tx.Save(dbUser).Scan(&dbUser); db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to update user with new password %s", db.Error.Error()))
		}

		return nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, "Failed to update user information, please try again later.", http.StatusInternalServerError)
		return
	}
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
