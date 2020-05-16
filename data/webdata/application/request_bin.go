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
	Method       string
	TimeReceived time.Time
	Proto        string
	RemoteAddr   string
	QueryStrings string
	Headers      string
	Forms        string
	Body         string
}
