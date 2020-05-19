package dnslookup

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//API dnslookup
type API struct{}

//DNSLookup look up dns
func (API) DNSLookup(c *gin.Context) {
	response := core.GetResponseInContext(c.Keys)
	request := &Request{}

	err := c.ShouldBind(&request)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid lookup request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	request.DomainName = strings.TrimSpace(request.DomainName)
	if !strings.HasPrefix(request.DomainName, "http") {
		request.DomainName = "http://" + request.DomainName
	}

	uri, err := url.Parse(request.DomainName)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid domain name",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if uri.Hostname() == "" {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Please enter a valid domain name",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	result := &Response{
		DomainName: uri.Hostname(),
		PTR:        make(map[string][]string),
	}

	ips, err := net.LookupIP(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Failed to lookup A records",
		})
	} else {
		for _, ip := range ips {
			result.IPAddress = append(result.IPAddress, ip.String())
			ptrs, err := net.LookupAddr(ip.String())
			if err != nil {
				utils.Logger.Error(err.Error())
				response.SetAlert(&data.AlertData{
					IsWarning: true,
					Message:   fmt.Sprintf("Failed to lookup PTR records for %s", ip.String()),
				})
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
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Failed to lookup CNAME records",
		})
	} else {
		for _, cname := range cnames {
			result.CNAME = append(result.CNAME, string(cname))
		}
	}

	nses, err := net.LookupNS(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Failed to lookup NS records",
		})
	} else {
		for _, ns := range nses {
			result.CNAME = append(result.NS, ns.Host)
		}
	}

	mxes, err := net.LookupMX(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Failed to lookup MX records",
		})
	} else {
		for _, mx := range mxes {
			result.CNAME = append(result.MX, mx.Host)
		}
	}

	txts, err := net.LookupTXT(uri.Hostname())
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Failed to lookup TXT records",
		})
	} else {
		for _, tx := range txts {
			result.CNAME = append(result.MX, tx)
		}
	}

	response.Data = result
	core.WriteResponse(c, 200, response)
}
