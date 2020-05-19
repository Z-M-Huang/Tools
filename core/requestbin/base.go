package requestbin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//PageData page data /app/request-bin
type PageData struct {
	ID              string
	URL             string
	VerificationKey string
	History         []*History
}

//History request history
type History struct {
	Method              string
	TimeReceived        time.Time
	Proto               string
	RemoteAddr          string
	QueryStrings        string
	Headers             map[string]string
	Cookies             []string
	Forms               map[string]string
	MultipartFormsFiles map[string]string
	Body                string
}

//CreateBinRequest /api/request-bin/create
type CreateBinRequest struct {
	IsPrivate bool `json:"isPrivate" xml:"isPrivate" form:"isPrivate"`
}

//CreateBinResponse response
type CreateBinResponse struct {
	ID              string
	VerificationKey string
}

//LoadRequestBinData load request bin data
func LoadRequestBinData(c *gin.Context) *PageData {
	id := c.Param("id")
	if id == "" {
		return nil
	}
	data := GetRequestBinHistory(id)
	if data != nil && data.VerificationKey != "" {
		val, err := c.Cookie("request_bin_verification_key")
		if err != nil || val == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/app/request-bin")
			c.Abort()
			return nil
		}
		if data.VerificationKey != val {
			c.Redirect(http.StatusTemporaryRedirect, "/app/request-bin")
			c.Abort()
			return nil
		}
	}
	return data
}
