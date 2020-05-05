package data

import "github.com/dgrijalva/jwt-go"

//JWTClaim web claim
type JWTClaim struct {
	ImageURL string `json:"image_url"`
	jwt.StandardClaims
}

//IsNil return true is claim is nil
func (c *JWTClaim) IsNil() bool {
	return c == nil
}
