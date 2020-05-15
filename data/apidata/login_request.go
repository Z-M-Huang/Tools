package apidata

//LoginRequest login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
