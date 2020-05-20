package dnslookup

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

var urlRe = regexp.MustCompile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

//API dnslookup
type API struct{}

func lookup(request *Request) (int, *data.APIResponse) {
	response := &data.APIResponse{}
	request.DomainName = strings.TrimSpace(request.DomainName)
	if !strings.HasPrefix(request.DomainName, "http") {
		request.DomainName = "http://" + request.DomainName
	}

	if !urlRe.Match([]byte(request.DomainName)) {
		response.Message = "Invalid domain name"
		return http.StatusBadRequest, response
	}

	uri, err := url.Parse(request.DomainName)
	if err != nil {
		return http.StatusInternalServerError, response
	}

	result := &Response{
		DomainName: uri.Hostname(),
		PTR:        make(map[string][]string),
	}

	ips, err := net.LookupIP(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Failed to lookup A records"
	} else {
		for _, ip := range ips {
			result.IPAddress = append(result.IPAddress, ip.String())
			ptrs, err := net.LookupAddr(ip.String())
			if err != nil {
				utils.Logger.Error(err.Error())
				response.Message = fmt.Sprintf("Failed to lookup PTR records for %s", ip.String())
			} else {
				for _, ptr := range ptrs {
					result.PTR[ip.String()] = append(result.PTR[ip.String()], ptr)
				}
			}
		}
	}

	cnames, err := net.LookupCNAME(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Failed to lookup CNAME records"
	} else {
		for _, cname := range cnames {
			result.CNAME = append(result.CNAME, string(cname))
		}
	}

	nses, err := net.LookupNS(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Failed to lookup NS records"
	} else {
		for _, ns := range nses {
			result.CNAME = append(result.NS, ns.Host)
		}
	}

	mxes, err := net.LookupMX(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Failed to lookup MX records"
	} else {
		for _, mx := range mxes {
			result.CNAME = append(result.MX, mx.Host)
		}
	}

	txts, err := net.LookupTXT(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Failed to lookup TXT records"
	} else {
		for _, tx := range txts {
			result.CNAME = append(result.MX, tx)
		}
	}
	response.Data = result
	return http.StatusOK, response
}

//DNSLookup look up dns
func (API) DNSLookup(c *gin.Context) {
	request := &Request{}

	err := c.ShouldBind(&request)
	if err != nil {
		core.WriteResponse(c, http.StatusBadRequest, &data.APIResponse{
			Message: "Bad Request",
		})
		return
	}

	status, response := lookup(request)

	core.WriteResponse(c, status, response)
}
