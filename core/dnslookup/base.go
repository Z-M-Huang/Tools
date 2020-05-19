package dnslookup

//Request dns lookup request
type Request struct {
	DomainName string `json:"domainName" xml:"domainName" form:"domainName" binding:"required"`
}

//Response dns lookup response
type Response struct {
	DomainName string
	IPAddress  []string
	CNAME      []string
	NS         []string
	MX         []string
	PTR        map[string][]string
	TXT        []string
}
