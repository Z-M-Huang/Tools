package home

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
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
		c.Set(utils.ClaimCtxKey, &account.JWTClaim{
			ImageURL: "https://localhost/favicon.ico",
			StandardClaims: jwt.StandardClaims{
				Id: "test@example.com",
			},
		})
	}
}

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	r.SetHTMLTemplate(loadTemplateRenderer())

	page := &Page{}
	r.GET("/", testHandler(), page.Home)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}
