package apidata

//UpdatePasswordRequest /api/account/update/password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
