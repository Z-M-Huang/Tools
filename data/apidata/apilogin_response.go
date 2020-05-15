package apidata

//APILoginResponse api login response
type APILoginResponse struct {
	AccessToken string `json:"access_token" xml:"access_token" form:"access_token" binding:"required"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
