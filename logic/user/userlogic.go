package userlogic

import (
	"fmt"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

//Find find db user
func Find(tx *gorm.DB, u *dbentity.User) error {
	if db := tx.Where(*u).First(&u); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save current user
func Save(tx *gorm.DB, u *dbentity.User) error {
	if db := tx.Save(u).Scan(&u); db.Error != nil {
		return db.Error
	}
	return nil
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
