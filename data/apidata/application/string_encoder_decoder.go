package application

//StringEncodeDecodeRequest request /api/string/encodedecode
type StringEncodeDecodeRequest struct {
	RequestString string `json:"requestString"`
	Type          string `json:"type"`
	Action        string `json:"action"`
}
