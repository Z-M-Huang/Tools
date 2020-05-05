package utils

import (
	"golang.org/x/crypto/bcrypt"
)

//HashAndSalt password
func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provied by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any alue you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		Logger.Error(err.Error())
		return ""
	}
	// GenerateFromPassword returns a byte slice s we need to
	// convert the byte to a string and return it
	return string(hash)
}

//ComparePasswords password
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting thehashed password from the DB it
	// will be a string so we'll need to convert it to a byt slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		Logger.Error(err.Error())
		return false
	}
	return true
}
