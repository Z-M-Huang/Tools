package iplocation

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/api/ip-location/get", api.Get)
	c.Request, _ = http.NewRequest("GET", "/api/ip-location/get", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRateLimit(t *testing.T) {
	assert.Empty(t, db.RedisSet(rateLimitKey, 5000, 24*time.Hour))
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.GET("/api/ip-location/get", api.Get)
	c.Request, _ = http.NewRequest("GET", "/api/ip-location/get", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
