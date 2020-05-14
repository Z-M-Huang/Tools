package apidata

//LoginRequest login request
type LoginRequest struct {
	Email    string `json:"email" xml:"email" form:"email" binding:"required"`
	Password string `json:"password" xml:"password" form:"password" binding:"required"`
}
