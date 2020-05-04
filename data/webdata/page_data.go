package webdata

//PageData page data contains all information for page to render
type PageData struct {
	LoginInfo   LoginData
	AlertInfo   AlertData
	ContentData interface{}
}

//AlertData used in header template
type AlertData struct {
	IsDanger  bool
	IsWarning bool
	Message   string
}

//LoginData page login info
type LoginData struct {
	Username string
	ImageURL string
}
