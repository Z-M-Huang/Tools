package apidata

//CreateAccountRequest /signup
type CreateAccountRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
