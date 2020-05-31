package shortlink

//Request to get short link
type Request struct {
	URL string `json:"url" xml:"url" form:"url" binding:"required"`
}
