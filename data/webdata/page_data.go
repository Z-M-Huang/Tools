package webdata

//PageData page data contains all information for page to render
type PageData struct {
	LoginInfo   LoginData
	ContentData interface{}
}

//LoginData page login info
type LoginData struct {
	Username string
	ImageURL string
}
