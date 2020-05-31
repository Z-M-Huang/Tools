package shortlink

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

func TestGetLink(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/shortlink/get", api.Get)
	form := url.Values{}
	form.Add("url", "https://example.com")
	c.Request, _ = http.NewRequest("POST", "/api/shortlink/get", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetLinkFailRequest(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/shortlink/get", api.Get)
	c.Request, _ = http.NewRequest("POST", "/api/shortlink/get", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectShortLink(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/s/:id", api.RedirectShortLink)
	c.Request, _ = http.NewRequest("GET", "/s/1", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
}

func TestRedirectShortLinkNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/s/:id", api.RedirectShortLink)
	c.Request, _ = http.NewRequest("GET", "/s/586715", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRedirectShortLinkIDNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/s/:id", api.RedirectShortLink)
	c.Request, _ = http.NewRequest("GET", "/s/", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRedirectShortLinkBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/s/:id", api.RedirectShortLink)
	c.Request, _ = http.NewRequest("GET", "/s/asdfasd", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedisRateLimit(t *testing.T) {
	db.RedisSet(getIPKey(""), 1, time.Hour)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/shortlink/get", api.Get)
	form := url.Values{}
	form.Add("url", "https://example.com")
	c.Request, _ = http.NewRequest("POST", "/api/shortlink/get", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	db.RedisDelete(getIPKey(""))
}

func TestRedisRateLimitFail(t *testing.T) {
	db.RedisSet(getIPKey(""), 10000, time.Hour)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/shortlink/get", api.Get)
	form := url.Values{}
	form.Add("url", "https://example.com")
	c.Request, _ = http.NewRequest("POST", "/api/shortlink/get", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	db.RedisDelete(getIPKey(""))
}
