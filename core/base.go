package core

import (
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//Response page and api response
type Response struct {
	Header *HeaderData `json:",omitempty"`
	Data   interface{}
}

//HeaderData header data
type HeaderData struct {
	Title           string
	Description     string
	ResourceVersion string
	PageStyle       *PageStyleData `json:",omitempty"`
	Nav             *NavData       `json:",omitempty"`
	Alert           *AlertData     `json:",omitempty"`
}

//AlertData used in web pages and api responses
type AlertData struct {
	IsInfo    bool
	IsSuccess bool
	IsWarning bool
	IsDanger  bool
	Message   string
}

//NavData nav bar
type NavData struct {
	StyleName string
	Login     *LoginData `json:",omitempty"`
}

//PageStyleData bootswatch styles
type PageStyleData struct {
	Name      string
	Link      string
	Integrity string
}

//LoginData page login info
type LoginData struct {
	Username string
	ImageURL string
}

//SetAlert set alert
func (r *Response) SetAlert(alert *AlertData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	r.Header.Alert = alert
}

//SetLogin set login
func (r *Response) SetLogin(login *LoginData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	if r.Header.Nav == nil {
		r.Header.Nav = &NavData{}
	}
	r.Header.Nav.Login = login
}

//SetNavStyleName set nav style name
func (r *Response) SetNavStyleName(style *PageStyleData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	r.Header.PageStyle = style
	if r.Header.Nav == nil {
		r.Header.Nav = &NavData{}
	}
	r.Header.Nav.StyleName = style.Name
}

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
func WriteUnexpectedError(c *gin.Context, response *Response) {
	response.SetAlert(&AlertData{
		IsDanger: true,
		Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
	})
	WriteResponse(c, http.StatusInternalServerError, response)
}

//GetResponseInContext get response struct from context
func GetResponseInContext(contextKey map[string]interface{}) *Response {
	response := contextKey[utils.ResponseCtxKey]
	if response == nil {
		return nil
	}
	return response.(*Response)
}

//SetCookie sets cookie
func SetCookie(c *gin.Context, cookieName, cookieVal string, expiresAt time.Time, httpOnly bool) {
	if data.Config.HTTPS {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, true, httpOnly)
	} else {
		c.SetCookie(cookieName, cookieVal, int(expiresAt.Sub(time.Now()).Seconds()), "/", data.Config.Host, false, httpOnly)
	}
}
