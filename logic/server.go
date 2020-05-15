package logic

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time, httpOnly bool) {
	if utils.Config.HTTPS {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", utils.Config.Host, true, httpOnly)
	} else {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", utils.Config.Host, false, httpOnly)
	}
}

//GetClaimFromCookieAndRenew get claim and renew
func GetClaimFromCookieAndRenew(c *gin.Context) (*data.JWTClaim, error) {
	val, err := c.Cookie(utils.SessionTokenKey)
	if err != nil || val == "" {
		return nil, nil
	}
	claim, err := isTokenValid(val)
	if err != nil {
		return nil, err
	}
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now().UTC()).Hours() < 24 {
		tokenStr, expiresAt, err := userlogic.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			SetCookie(c, utils.SessionTokenKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

//GetClaimFromHeaderAndRenew get claim and renew
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
		tokenStr, expiresAt, err := userlogic.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			SetCookie(c, utils.SessionTokenKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

func isTokenValid(token string) (*data.JWTClaim, error) {
	claims := &data.JWTClaim{}
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
