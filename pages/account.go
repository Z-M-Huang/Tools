package pages

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	tokenStr, expireAt, err := utils.GenerateJWTToken("Google", user.Email, user.Name, user.Picture)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		cookie := &http.Cookie{
			Name:       utils.SessionTokenKey,
			Value:      tokenStr,
			Path:       "/",
			Domain:     utils.Config.Host,
			Expires:    expireAt,
			RawExpires: expireAt.String(),
		}
		if utils.Config.IsDebug {
			cookie.Domain = "localhost"
		}
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

//LoginPage /login
func LoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.Templates.ExecuteTemplate(w, "login.gohtml", &data.Response{})
}

//AccountPage /account requires claim
func AccountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageData := &data.Response{}
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	pageData.Login = data.LoginData{
		Username: claim.Subject,
		ImageURL: claim.ImageURL,
	}
	user, err := GetUserInfoFromDB(claim.Id)
	if err != nil {
		pageData.Alert.IsDanger = true
		pageData.Alert.Message = err.Error()
	} else {
		pageData.Data = user
	}
	utils.Templates.ExecuteTemplate(w, "account.gohtml", pageData)
}

//GetUserInfoFromDB get user info from database
func GetUserInfoFromDB(emailAddress string) (*dbentity.User, error) {
	user := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: emailAddress,
	}).First(&user); db.RecordNotFound() {
		return nil, fmt.Errorf("User not found associated with email: %s", emailAddress)
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		return nil, errors.New("Internal Error: failed to get user info, please try again later")
	}
	return user, nil
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
