package apidata

//CreateAccountRequest /signup
type CreateAccountRequest struct {
	Email           string `json:"email" xml:"email" form:"email" binding:"required"`
	Username        string `json:"username" xml:"username" form:"username" binding:"required"`
	Password        string `json:"password" xml:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" xml:"confirmPassword" form:"confirmPassword" binding:"required"`
}
