package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data/apidata"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config
var authTokenKey = "session_token"

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
	tokenStr, expireAt, err := generateJWTToken("Google", user.Email, user.Name, user.Picture)
	if err != nil {
		utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
	} else {
		cookie := &http.Cookie{
			Name:       authTokenKey,
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

//GetUserInfo get user info for page and renew token if needed
func GetUserInfo(w http.ResponseWriter, r *http.Request) (*webdata.LoginData, error) {
	cookies := r.Cookies()
	for _, c := range cookies {
		fmt.Println(fmt.Sprintf("%v, %v", c.Name, c.Value))
	}
	cookie, err := r.Cookie(authTokenKey)
	if err != nil {
		return nil, nil
	}
	claim, err := isTokenValid(cookie.Value)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().Sub(time.Unix(claim.ExpiresAt, 0)).Hours() < 24 {
		tokenStr, expireAt, err := generateJWTToken("Google", claim.Id, claim.Subject, claim.ImageURL)
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
	}

	return &webdata.LoginData{
		Username: claim.Subject,
		ImageURL: claim.ImageURL,
	}, nil
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

func generateJWTToken(audience, emailAddress, username, imageURL string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claim := &webdata.JWTClaim{
		ImageURL: imageURL,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Id:        emailAddress,
			Audience:  audience,
			Subject:   username,
			Issuer:    utils.Config.Host,
		},
	}
	claim.ImageURL = imageURL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(utils.Config.JwtKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("failed to generate token %s", err.Error())
	}
	return tokenStr, expiresAt, nil
}

func isTokenValid(token string) (*webdata.JWTClaim, error) {
	claims := &webdata.JWTClaim{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.Config.JwtKey, nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Unauthenticated")
	}

	if !tkn.Valid || !claims.VerifyIssuer(utils.Config.Host, true) {
		return nil, fmt.Errorf("Invalid Token")
	}

	return claims, nil
}
