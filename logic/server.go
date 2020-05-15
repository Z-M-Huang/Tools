package logic

import (
	"io"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/utils"
)

// GzipResponseWriter is a Struct for manipulating io writer
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (res GzipResponseWriter) Write(b []byte) (int, error) {
	if "" == res.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		res.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return res.Writer.Write(b)
}

//SetCookie sets cookie
func SetCookie(w http.ResponseWriter, cookieName, cookieVal string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:       cookieName,
		Value:      cookieVal,
		Path:       "/",
		Domain:     utils.Config.Host,
		Expires:    expiresAt,
		RawExpires: expiresAt.String(),
	})
}
