package application

//CreateBinRequest /api/request-bin/create
type CreateBinRequest struct {
	IsPrivate bool `json:"isPrivate" xml:"isPrivate" form:"isPrivate"`
}

//CreateBinResponse response
type CreateBinResponse struct {
	URL             string
	VerificationKey string
}
