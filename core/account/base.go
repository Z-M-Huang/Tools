package account

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//JWTClaim web claim
type JWTClaim struct {
	ImageURL string `json:"image_url"`
	jwt.StandardClaims
}

//CreateAccountRequest /signup
type CreateAccountRequest struct {
	Email           string `json:"email" xml:"email" form:"email" binding:"required"`
	Username        string `json:"username" xml:"username" form:"username" binding:"required"`
	Password        string `json:"password" xml:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" xml:"confirmPassword" form:"confirmPassword" binding:"required"`
}

//UpdatePasswordRequest /api/account/update/password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" xml:"currentPassword" form:"currentPassword" binding:"required"`
	Password        string `json:"password" xml:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" xml:"confirmPassword" form:"confirmPassword" binding:"required"`
}

//GoogleUserInfo user info
type GoogleUserInfo struct {
	ID            string `json:"id"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Local         string `json:"local"`
	Email         string `json:"Email"`
	GivenName     string `json:"GivenName"`
	VerifiedEmail bool   `json:"verified_email"`
}

//LoginRequest login request
type LoginRequest struct {
	Email    string `json:"email" xml:"email" form:"email" binding:"required"`
	Password string `json:"password" xml:"password" form:"password" binding:"required"`
}

//LoginResponse login response
type LoginResponse struct {
	IsSuccess bool
	Redirect  string
}

//PageData /account
type PageData struct {
	HasPassword bool
}

//GetClaimInContext get claim struct from context
func GetClaimInContext(contextKey map[string]interface{}) *JWTClaim {
	claim := contextKey[utils.ClaimCtxKey]
	if claim == nil {
		return nil
	}
	return claim.(*JWTClaim)
}

//GetClaimFromCookieAndRenew get claim and renew
func GetClaimFromCookieAndRenew(c *gin.Context) (*JWTClaim, error) {
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
			core.SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

//GetClaimFromHeaderAndRenew get claim and renew. Since auth token is httponly, it will not really be able to get from javascript
func GetClaimFromHeaderAndRenew(c *gin.Context) (*JWTClaim, error) {
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
			core.SetCookie(c, utils.SessionCookieKey, tokenStr, expiresAt, false)
		}
	}
	return claim, nil
}

//GenerateJWTToken generates JWT token
func GenerateJWTToken(audience, emailAddress, username, imageURL string) (string, time.Time, error) {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	claim := &JWTClaim{
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

func isTokenValid(token string) (*JWTClaim, error) {
	claims := &JWTClaim{}
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
