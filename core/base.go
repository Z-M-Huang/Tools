package core

import (
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//WriteResponse Write api response
func WriteResponse(c *gin.Context, code int, obj interface{}) {
	switch c.NegotiateFormat(gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.JSON(code, obj)
	case gin.MIMEXML:
		c.XML(code, obj)
	case gin.MIMEYAML:
		c.YAML(code, obj)
	default:
		c.String(http.StatusNotAcceptable, "Not Acceptable")
	}
}

//WriteUnexpectedError Write unexpected api response
func WriteUnexpectedError(c *gin.Context, response *data.PageResponse) {
	response.SetAlert(&data.AlertData{
		IsDanger: true,
		Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
	})
	WriteResponse(c, http.StatusInternalServerError, response)
}

//GetResponseInContext get response struct from context
func GetResponseInContext(contextKey map[string]interface{}) *data.PageResponse {
	response := contextKey[utils.ResponseCtxKey]
	if response == nil {
		return nil
	}
	return response.(*data.PageResponse)
}

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time, httpOnly bool) {
	if data.Config.HTTPS {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, true, httpOnly)
	} else {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, false, httpOnly)
	}
}
