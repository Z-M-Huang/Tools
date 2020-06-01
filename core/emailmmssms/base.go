package emailmmssms

//Request for /api/email-mms-sms/send
type Request struct {
	ToNumber string `json:"toNumber" xml:"toNumber" form:"toNumber" binding:"required"`
	Carrier  string `json:"carrier" xml:"carrier" form:"carrier" binding:"required"`
	Subject  string `json:"subject" xml:"subject" form:"subject"`
	Content  string `json:"content" xml:"content" form:"content"`
}

//LookupRequest for /api/email-mms-sms/lookup
type LookupRequest struct {
	Number string `json:"phone_number" xml:"phone_number" form:"phone_number" binding:"required"`
}

//LookupResponse for /api/email-mms-sms/lookup
type LookupResponse struct {
	CountryCode    int64   `json:"country_code" xml:"country_code"`
	CountryCodeISO string  `json:"country_code_iso" xml:"country_code_iso"`
	Location       string  `json:"location" xml:"location"`
	Latitude       float64 `json:"location_latitude" xml:"location_latitude"`
	Longitude      float64 `json:"location_longitude" xml:"location_longitude"`
	NationalNumber int64   `json:"national_number" xml:"national_number"`
	Type           string  `json:"number_type" xml:"number_type"`
	IsValid        bool    `json:"is_valid_number" xml:"is_valid_number"`
	Carrier        string  `json:"carrier" xml:"carrier"`
	E164Formatted  string  `json:"phone_number_e164" xml:"phone_number_e164"`
}
