package emailmmssms

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

func TestSendFailwithNoConfig(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmailRateLimit(t *testing.T) {
	db.RedisSet(getTotalEmailKey(), 10, 2*time.Second)
	maxDailyEmailAmount = 15
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	db.RedisDelete(getTotalEmailKey())
}

func TestIPLimit(t *testing.T) {
	db.RedisSet(getRedisIPKey(""), 2000, 2*time.Second)
	db.RedisDelete(getRedisToNumberKey("1234567890"))
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestToNumberLimit(t *testing.T) {
	db.RedisSet(getRedisToNumberKey("1234567890"), 2000, 2*time.Second)
	db.RedisDelete(getRedisIPKey(""))
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestInvalidNumber(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "123456789")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNoContentNumber(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&T")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnknownCarrier(t *testing.T) {
	w := httptest.NewRecorder()
	maxDailyEmailAmount = 15000
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/api/email-mms-sms/send", api.Send)
	form := url.Values{}
	form.Add("content", "test")
	form.Add("subject", "test")
	form.Add("toNumber", "1234567890")
	form.Add("carrier", "AT&TTTT")
	c.Request, _ = http.NewRequest("POST", "/api/email-mms-sms/send", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
