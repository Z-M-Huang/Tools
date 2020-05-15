package application

//DNSLookupRequest dns lookup request
type DNSLookupRequest struct {
	DomainName string `json:"domainName"`
}

//DNSLookupResponse dns lookup response
type DNSLookupResponse struct {
	DomainName string
	IPAddress  []string
	CNAME      []string
	NS         []string
	MX         []string
	PTR        map[string][]string
	TXT        []string
}
