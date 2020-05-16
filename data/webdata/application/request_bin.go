package application

import "time"

//RequestBinPageData page data /app/request-bin
type RequestBinPageData struct {
	ID              string
	URL             string
	VerificationKey string
	History         []*RequestHistory
}

//RequestHistory request history
type RequestHistory struct {
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
