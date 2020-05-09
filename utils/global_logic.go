package utils

import (
	"net/http"
	"time"
)

//SetCookie sets cookie
func SetCookie(w http.ResponseWriter, cookieName, cookieVal string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:       cookieName,
		Value:      cookieVal,
		Path:       "/",
		Domain:     Config.Host,
		Expires:    expiresAt,
		RawExpires: expiresAt.String(),
	})
}
