package portchecker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckPort(t *testing.T) {
	open := checkPort("google.com", "80", "tcp")
	assert.True(t, open)
}

func TestCheckFail(t *testing.T) {
	open := checkPort("googlegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegooglegoogle.com", "80", "tcp")
	assert.False(t, open)
}

func TestCheck(t *testing.T) {
	db.RedisDelete("APP_PORT_CHECKER_")
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "80")
	form.Add("type", "tcp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCheckBindFail(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "asdfasdf")
	form.Add("type", "tcp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCheckPortNumberTooBigFail(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "65540")
	form.Add("type", "tcp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCheckPortNumberTooSmallFail(t *testing.T) {
	db.RedisDelete("APP_PORT_CHECKER_")
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "-1")
	form.Add("type", "tcp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCheckPortTypeFail(t *testing.T) {
	db.RedisDelete("APP_PORT_CHECKER_")
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "80")
	form.Add("type", "icmp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCheckTooManyRequestsFail(t *testing.T) {
	db.RedisSet("APP_PORT_CHECKER_", 1, 3*time.Second)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/portchecker/check", page.Check)
	form := url.Values{}
	form.Add("host", "google.com")
	form.Add("port", "80")
	form.Add("type", "tcp")
	c.Request, _ = http.NewRequest("POST", "/api/portchecker/check", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.NotEmpty(t, w.Body)
}
