package apidata

//UpdatePasswordRequest /api/account/update/password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" xml:"currentPassword" form:"currentPassword" binding:"required"`
	Password        string `json:"password" xml:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" xml:"confirmPassword" form:"confirmPassword" binding:"required"`
}
