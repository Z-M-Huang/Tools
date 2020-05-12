package data

//Response page and api response
type Response struct {
	Header *HeaderData `json:",omitempty"`
	Data   interface{}
}

//HeaderData header data
type HeaderData struct {
	Title           string
	ResourceVersion string
	Login           *LoginData `json:",omitempty"`
	Alert           *AlertData
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

//SetAlert set alert
func (r *Response) SetAlert(alert *AlertData) {
	if r.Header == nil {
		r.Header = &HeaderData{
			Alert: alert,
		}
	} else {
		r.Header.Alert = alert
	}
}

//SetLogin set login
func (r *Response) SetLogin(login *LoginData) {
	if r.Header == nil {
		r.Header = &HeaderData{
			Login: login,
		}
	} else {
		r.Header.Login = login
	}
}
