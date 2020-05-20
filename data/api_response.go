package data

//APIResponse api response
type APIResponse struct {
	ErrorMessage string      `json:",omitempty"`
	Data         interface{} `json:",omitempty"`
}
