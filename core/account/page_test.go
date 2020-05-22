package account

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testNoClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(utils.ResponseCtxKey, &data.PageResponse{
			Header: &data.HeaderData{
				ResourceVersion: "test",
				PageStyle: &data.PageStyleData{
					Name:      "Default",
					Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
					Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk",
				},
				Nav: &data.NavData{
					StyleName: "Default",
				},
			},
		})
	}
}

func testClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &CreateAccountRequest{
			Email:           "TestAccountPageWithUser@example.com",
			Username:        "TestAccountPageWithUser",
			Password:        "abcdef123456",
			ConfirmPassword: "abcdef123456",
		}
		signUp(request)

		c.Set(utils.ResponseCtxKey, &data.PageResponse{
			Header: &data.HeaderData{
				ResourceVersion: "test",
				PageStyle: &data.PageStyleData{
					Name:      "Default",
					Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
					Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk",
				},
				Nav: &data.NavData{
					StyleName: "Default",
				},
			},
		})
		c.Set(utils.ClaimCtxKey, &JWTClaim{
			ImageURL: "https://localhost/favicon.ico",
			StandardClaims: jwt.StandardClaims{
				Id: request.Email,
			},
		})
	}
}

func testInvalidClaimHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(utils.ResponseCtxKey, &data.PageResponse{
			Header: &data.HeaderData{
				ResourceVersion: "test",
				PageStyle: &data.PageStyleData{
					Name:      "Default",
					Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
					Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk",
				},
				Nav: &data.NavData{
					StyleName: "Default",
				},
			},
		})
		c.Set(utils.ClaimCtxKey, &JWTClaim{
			ImageURL: "https://localhost/favicon.ico",
			StandardClaims: jwt.StandardClaims{
				Id: "I have no idea where this user is",
			},
		})
	}
}

func TestSignupPageNoAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/signup", testNoClaimHandler(), page.Signup)
	c.Request, _ = http.NewRequest("GET", "/signup", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignupPageAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/signup", testClaimHandler(), page.Signup)
	c.Request, _ = http.NewRequest("GET", "/signup", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestLoginPageNoAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/login", testNoClaimHandler(), page.Login)
	c.Request, _ = http.NewRequest("GET", "/login?redirect=/account", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginPageAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/login", testClaimHandler(), page.Login)
	c.Request, _ = http.NewRequest("GET", "/login", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAccountPageNoAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/account", testNoClaimHandler(), page.Account)
	c.Request, _ = http.NewRequest("GET", "/account", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestAccountPageAuth(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/account", testClaimHandler(), page.Account)
	c.Request, _ = http.NewRequest("GET", "/account", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccountPageInvalidClaim(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	tplt := loadTemplateRenderer()
	r.SetHTMLTemplate(tplt)

	page := &Page{}
	r.GET("/account", testInvalidClaimHandler(), page.Account)
	c.Request, _ = http.NewRequest("GET", "/account", nil)

	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}
