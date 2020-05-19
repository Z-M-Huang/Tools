package stringencoderdecoder

//Request request /api/string/encodedecode
type Request struct {
	RequestString string `json:"requestString" xml:"requestString" form:"requestString" binding:"required"`
	Type          string `json:"type" xml:"type" form:"type" binding:"required"`
	Action        string `json:"action" xml:"action" form:"action" binding:"required"`
}
