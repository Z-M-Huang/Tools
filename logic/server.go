package logic

import (
	"time"

	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time) {
	c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", utils.Config.Host, true, false)
}
