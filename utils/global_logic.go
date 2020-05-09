package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/dgrijalva/jwt-go"
)

//GenerateJWTToken generates JWT token
func GenerateJWTToken(audience, emailAddress, username, imageURL string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claim := &data.JWTClaim{
		ImageURL: imageURL,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Id:        emailAddress,
			Audience:  audience,
			Subject:   username,
			Issuer:    Config.Host,
		},
	}
	claim.ImageURL = imageURL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(Config.JwtKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("failed to generate token %s", err.Error())
	}
	return tokenStr, expiresAt, nil
}

//GetApplicationUsed saved in cookie
func GetApplicationUsed(r *http.Request) ([]string, error) {
	var usedApps []string
	usedAppCookie, err := r.Cookie(UsedTokenKey)
	if err == http.ErrNoCookie {
		return usedApps, nil
	} else if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(usedAppCookie.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(decoded, &usedApps)
	if err != nil {
		return nil, err
	}
	return usedApps, nil
}

//GetApplicationsByName get application by name
func GetApplicationsByName(name string) *webdata.AppCard {
	for _, category := range AppList {
		for _, app := range category.AppCards {
			if app.Name == name {
				return app
			}
		}
	}
	return nil
}

//SetCookie sets cookie
func SetCookie(w http.ResponseWriter, cookieName, cookieVal string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:       cookieName,
		Value:      cookieVal,
		Path:       "/",
		Domain:     Config.Host,
		Expires:    expiresAt,
		RawExpires: expiresAt.String(),
	})
}
