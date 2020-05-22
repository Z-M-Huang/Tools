package application

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testHandler() gin.HandlerFunc {
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

func TestRenderApplicationPage(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)
			r.SetHTMLTemplate(loadTemplateRenderer())

			page := &Page{}
			r.GET("/app/:name", testHandler(), page.RenderApplicationPage)
			c.Request, _ = http.NewRequest("GET", "/app/"+app.Name, nil)
			c.Request.AddCookie(&http.Cookie{
				Name:  utils.UsedTokenKey,
				Value: "",
			})
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestRenderApplicationPageWithCookie(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)
			r.SetHTMLTemplate(loadTemplateRenderer())

			page := &Page{}
			r.GET("/app/:name", testHandler(), page.RenderApplicationPage)
			c.Request, _ = http.NewRequest("GET", "/app/"+app.Name, nil)
			c.Request.AddCookie(&http.Cookie{
				Name:  utils.UsedTokenKey,
				Value: "WyJIaUxvIFNpbXVsYXRvciIsIkVuY29kZXIgRGVjb2RlciIsIlFSIENvZGUiLCJETlMgTG9va3VwIiwiUmVxdWVzdCBCaW4iLCJLZWxseSBDcml0ZXJpb24iXQ%3D%3D",
			})
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestRenderApplicationPageWithInvalidCookie(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)
			r.SetHTMLTemplate(loadTemplateRenderer())

			page := &Page{}
			r.GET("/app/:name", testHandler(), page.RenderApplicationPage)
			c.Request, _ = http.NewRequest("GET", "/app/"+app.Name, nil)
			c.Request.AddCookie(&http.Cookie{
				Name:  utils.UsedTokenKey,
				Value: "a",
			})
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestRenderApplicationPageWithNoCookie(t *testing.T) {
	for _, category := range GetAppList() {
		for _, app := range category.AppCards {
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, r := gin.CreateTestContext(w)
			r.SetHTMLTemplate(loadTemplateRenderer())

			page := &Page{}
			r.GET("/app/:name", testHandler(), page.RenderApplicationPage)
			c.Request, _ = http.NewRequest("GET", "/app/"+app.Name, nil)
			r.ServeHTTP(w, c.Request)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestRenderApplicationPageWithInvalidApp(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	r.SetHTMLTemplate(loadTemplateRenderer())

	page := &Page{}
	r.GET("/app/:name", testHandler(), page.RenderApplicationPage)
	c.Request, _ = http.NewRequest("GET", "/app/", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusNotFound, w.Code)
}