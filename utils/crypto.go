package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

//Encrypt encrypts text with key
func Encrypt(text, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

//Decrypt decrypts text with key
func Decrypt(text, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

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
