package data

type Response struct {
	Login LoginData
	Alert AlertData
	Data  interface{}
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
