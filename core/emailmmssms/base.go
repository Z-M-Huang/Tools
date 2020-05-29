package emailmmssms

import "time"

//Request for /api/email-mms-sms/send
type Request struct {
	ToNumber string `json:"toNumber" xml:"toNumber" form:"toNumber"`
	Carrier  string `json:"carrier" xml:"carrier" form:"carrier"`
	Subject  string `json:"subject" xml:"subject" form:"subject"`
	Content  string `json:"content" xml:"content" form:"content"`
}

//MessageHistory replied messages
type MessageHistory struct {
	DateReceived time.Time
	Content      string
}
