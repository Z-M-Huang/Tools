package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data/apidata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config
var authTokenKey = "session_token"

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/google_oauth", utils.Config.Host),
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
	tokenStr, expireAt, err := generateJWTToken("Google", user.Email, user.GivenName)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:       authTokenKey,
			Value:      tokenStr,
			Path:       "/",
			Domain:     utils.Config.Host,
			Expires:    expireAt,
			RawExpires: expireAt.String(),
		})
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

func generateJWTToken(audience, emailAddress, username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claim := &jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Id:        emailAddress,
		Audience:  audience,
		Subject:   username,
		Issuer:    utils.Config.Host,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(utils.Config.JwtKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("failed to generate token %s", err.Error())
	}
	return tokenStr, expiresAt, nil
}
