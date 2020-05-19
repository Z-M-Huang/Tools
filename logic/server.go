package logic

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time, httpOnly bool) {
	if data.Config.HTTPS {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, true, httpOnly)
	} else {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, false, httpOnly)
	}
}

//GetClaimFromCookieAndRenew get claim and renew
func GetClaimFromCookieAndRenew(c *gin.Context) (*data.JWTClaim, error) {
	val, err := c.Cookie(utils.SessionCookieKey)
	if err != nil || val == "" {
		return nil, nil
	}
	claim, err := isTokenValid(val)
	if err != nil {
		return nil, err
	}
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now().UTC()).Hours() < 24 {
		tokenStr, expiresAt, err := GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

//GetClaimFromHeaderAndRenew get claim and renew. Since auth token is httponly, it will not really be able to get from javascript
func GetClaimFromHeaderAndRenew(c *gin.Context) (*data.JWTClaim, error) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.Contains(token, "Bearer ") {
		return nil, errors.New("Unauthorized")
	}

	token = strings.ReplaceAll(token, "Bearer ", "")
	claim, err := isTokenValid(token)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now().UTC()).Hours() < 24 {
		tokenStr, expiresAt, err := GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

//GenerateJWTToken generates JWT token
func GenerateJWTToken(audience, emailAddress, username, imageURL string) (string, time.Time, error) {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	claim := &data.JWTClaim{
		ImageURL: imageURL,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Id:        emailAddress,
			Audience:  audience,
			Subject:   username,
			Issuer:    data.Config.Host,
		},
	}
	claim.ImageURL = imageURL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(data.Config.JwtKey)
	if err != nil {
		return "", time.Now(), fmt.Errorf("failed to generate token %s", err.Error())
	}
	return tokenStr, expiresAt, nil
}

func isTokenValid(token string) (*data.JWTClaim, error) {
	claims := &data.JWTClaim{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return data.Config.JwtKey, nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Unauthenticated")
	}

	if !tkn.Valid || !claims.VerifyIssuer(data.Config.Host, true) {
		return nil, fmt.Errorf("Invalid Token")
	}

	return claims, nil
}
