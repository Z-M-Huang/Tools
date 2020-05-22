package data

//PageResponse page response
type PageResponse struct {
	Header *HeaderData `json:",omitempty"`
	Data   interface{}
}

//HeaderData header data
type HeaderData struct {
	Title           string
	Description     string
	ResourceVersion string
	PageStyle       *PageStyleData `json:",omitempty"`
	Nav             *NavData       `json:",omitempty"`
	Alert           *AlertData     `json:",omitempty"`
}

//AlertData used in web pages and api responses
type AlertData struct {
	IsInfo    bool
	IsSuccess bool
	IsWarning bool
	IsDanger  bool
	Message   string
}

//NavData nav bar
type NavData struct {
	StyleName string
	Login     *LoginData `json:",omitempty"`
}

//PageStyleData bootswatch styles
type PageStyleData struct {
	Name      string
	Link      string
	Integrity string
}

//LoginData page login info
type LoginData struct {
	Username string
	ImageURL string
}

//SetAlert set alert
func (r *PageResponse) SetAlert(alert *AlertData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	r.Header.Alert = alert
}

//SetLogin set login
func (r *PageResponse) SetLogin(login *LoginData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	if r.Header.Nav == nil {
		r.Header.Nav = &NavData{}
	}
	r.Header.Nav.Login = login
}

//SetNavStyleName set nav style name
func (r *PageResponse) SetNavStyleName(style *PageStyleData) {
	if r.Header == nil {
		r.Header = &HeaderData{}
	}
	r.Header.PageStyle = style
	if r.Header.Nav == nil {
		r.Header.Nav = &NavData{}
	}
	r.Header.Nav.StyleName = style.Name
}
