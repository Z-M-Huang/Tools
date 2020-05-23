package qrcode

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateQRCode(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "H")
	form.Add("backColor", "ffffff")
	form.Add("foreColor", "000000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeLowLevel(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "L")
	form.Add("backColor", "ffffff")
	form.Add("foreColor", "000000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeMediumLevel(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "M")
	form.Add("backColor", "ffffff")
	form.Add("foreColor", "000000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeQuartileLevel(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "Q")
	form.Add("backColor", "ffffff")
	form.Add("foreColor", "000000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeNoSize(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeSizeTooLarge(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "5000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeSizeTooSmall(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "-1")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("size", "256")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeInvalidLevel(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "X")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeBadBackColor(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "L")
	form.Add("backColor", "xxxxxx")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeBadForeColor(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "L")
	form.Add("foreColor", "xxxxxx")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateQRCodeToomanyRequest(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	redisKey := getRedisKey("")
	db.RedisSet(redisKey, -1, 24*time.Hour)

	page := &API{}
	r.POST("/api/qr-code/create", page.CreateQRCode)
	form := url.Values{}
	form.Add("content", "test content")
	form.Add("size", "256")
	form.Add("level", "H")
	form.Add("backColor", "ffffff")
	form.Add("foreColor", "000000")
	c.Request, _ = http.NewRequest("POST", "/api/qr-code/create", strings.NewReader(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.NotEmpty(t, w.Body)
}
