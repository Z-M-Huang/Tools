package data

//APIResponse api response
type APIResponse struct {
	Message string      `json:",omitempty"`
	Data    interface{} `json:",omitempty"`
}
