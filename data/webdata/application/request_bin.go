package application

import "net/http"

//RequestBinPageData page data /app/request-bin
type RequestBinPageData struct {
	ID              string
	URL             string
	VerificationKey string
	History         []*http.Request
}
