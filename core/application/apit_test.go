package application

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &account.CreateAccountRequest{
			Email:           "test@example.com",
			Username:        "testUsername",
			Password:        "abcdef123456",
			ConfirmPassword: "abcdef123456",
		}
		account.SignUp(request)
		c.Set(utils.ClaimCtxKey, &account.JWTClaim{
			ImageURL: "https://localhost/favicon.ico",
			StandardClaims: jwt.StandardClaims{
				Id: request.Email,
			}})
		c.Next()
	}
}

func testUnkownClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(utils.ClaimCtxKey, &account.JWTClaim{
			ImageURL: "https://localhost/favicon.ico",
			StandardClaims: jwt.StandardClaims{
				Id: "unknwon@example.com",
			}})
		c.Next()
	}
}

func TestLike(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)

			api := &API{}
			r.POST("/app/:name/like", testHandler(), testClaimHandler(), api.Like)
			c.Request, _ = http.NewRequest("POST", "/app/"+app.Name+"/like", nil)
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestLikeFailNoAppName(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/app/:name/like", testHandler(), testClaimHandler(), api.Like)
	c.Request, _ = http.NewRequest("POST", "/app//like", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLikeFailUnknownAppName(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/app/:name/like", testHandler(), testClaimHandler(), api.Like)
	c.Request, _ = http.NewRequest("POST", "/app/ohmy/like", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLikeFailUnknownUser(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)

			api := &API{}
			r.POST("/app/:name/like", testHandler(), testUnkownClaimHandler(), api.Like)
			c.Request, _ = http.NewRequest("POST", "/app/"+app.Name+"/like", nil)
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	}
}

func TestDisLike(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)

			api := &API{}
			r.POST("/app/:name/dislike", testHandler(), testClaimHandler(), api.Dislike)
			c.Request, _ = http.NewRequest("POST", "/app/"+app.Name+"/dislike", nil)
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestDisLikeFailNoAppName(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/app/:name/dislike", testHandler(), testClaimHandler(), api.Dislike)
	c.Request, _ = http.NewRequest("POST", "/app//dislike", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDisLikeFailUnknownAppName(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/app/:name/dislike", testHandler(), testClaimHandler(), api.Dislike)
	c.Request, _ = http.NewRequest("POST", "/app/ohmy/dislike", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDisLikeFailUnknownUser(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)

			api := &API{}
			r.POST("/app/:name/dislike", testHandler(), testUnkownClaimHandler(), api.Dislike)
			c.Request, _ = http.NewRequest("POST", "/app/"+app.Name+"/dislike", nil)
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		}
	}
}
