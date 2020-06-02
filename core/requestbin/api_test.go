package requestbin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	bin := create(false)

	assert.NotEmpty(t, bin)
	assert.NotEmpty(t, bin.ID)
	assert.NotEmpty(t, bin.URL)
	assert.Empty(t, bin.VerificationKey)

	data.Config.HTTPS = true
	privateBin := create(true)
	assert.NotEmpty(t, privateBin)
	assert.NotEmpty(t, privateBin.ID)
	assert.NotEmpty(t, privateBin.URL)
	assert.NotEmpty(t, privateBin.VerificationKey)
}

func TestCreateRequestBin(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)

	api := &API{}
	r.POST("/request-bin/create", api.CreateRequestBin)
	form := url.Values{}
	form.Add("isPrivate", "true")
	c.Request, _ = http.NewRequest("POST", "/request-bin/create", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRequestBinFail(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	api := &API{}
	r.POST("/request-bin/create", api.CreateRequestBin)
	c.Request, _ = http.NewRequest("POST", "/request-bin/create", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRequestBinHistory(t *testing.T) {
	binData := create(false)

	redisBinData := GetRequestBinHistory(binData.ID)
	assert.NotEmpty(t, redisBinData)
}

func TestGetRequestBinHistoryFail(t *testing.T) {
	redisBinData := GetRequestBinHistory("123456")
	assert.Empty(t, redisBinData)
}

func TestRequestIn(t *testing.T) {
	bin := create(false)
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch}
	for _, m := range methods {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(w)
		api := &API{}
		form := url.Values{}
		form.Add("test", "test")
		r.Any("/api/request-bin/receive/:id", api.RequestIn)
		c.Request, _ = http.NewRequest(m, "/api/request-bin/receive/"+bin.ID, bytes.NewBufferString(form.Encode()))
		c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, c.Request)
		assert.Equal(t, http.StatusOK, w.Code, m)
	}
}

func TestRequestInFailNoID(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	api := &API{}
	form := url.Values{}
	form.Add("test", "test")
	r.Any("/api/request-bin/receive/:id", api.RequestIn)
	c.Request, _ = http.NewRequest("POST", "/api/request-bin/receive/", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRequestInFailNotFoundID(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	api := &API{}
	form := url.Values{}
	form.Add("test", "test")
	r.Any("/api/request-bin/receive/:id", api.RequestIn)
	c.Request, _ = http.NewRequest("POST", "/api/request-bin/receive/123", bytes.NewBufferString(form.Encode()))
	c.Request.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
