package data

//Response page and api response
type Response struct {
	Header *HeaderData `json:",omitempty"`
	Alert  AlertData
	Data   interface{}
}

//HeaderData header data
type HeaderData struct {
	Title           string
	ResourceVersion string
	Login           *LoginData `json:",omitempty"`
}

//LoginData page login info
type LoginData struct {
	Username string
	ImageURL string
}

//AlertData used in web pages and api responses
type AlertData struct {
	IsInfo    bool
	IsSuccess bool
	IsWarning bool
	IsDanger  bool
	Message   string
}
