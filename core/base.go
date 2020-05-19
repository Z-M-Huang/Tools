package core

import (
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//GetResponseInContext get response struct from context
func GetResponseInContext(contextKey map[string]interface{}) *data.Response {
	response := contextKey[utils.ResponseCtxKey]
	if response == nil {
		return nil
	}
	return response.(*data.Response)
}

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time, httpOnly bool) {
	if data.Config.HTTPS {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, true, httpOnly)
	} else {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, false, httpOnly)
	}
}
