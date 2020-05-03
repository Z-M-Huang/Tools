package apidata

//GoogleUserInfo user info
type GoogleUserInfo struct {
	ID string `json:"id"`
	FamilyName string `json:"family_name"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	Local string `json:"local"`
	Email string `json:"Email"`
	GivenName string `json:"GivenName"`
	VerifiedEmail bool `json:"verified_email"`
}