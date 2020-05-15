package api

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
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
func WriteUnexpectedError(c *gin.Context, response *data.Response) {
	response.SetAlert(&data.AlertData{
		IsDanger: true,
		Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
	})
	WriteResponse(c, http.StatusInternalServerError, response)
}
