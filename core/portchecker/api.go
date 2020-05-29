package portchecker

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/gin-gonic/gin"
)

//API Port Checker API
type API struct{}

func checkPort(host, port, portType string) bool {
	timeout := 3 * time.Second
	var conn net.Conn
	conn, err := net.DialTimeout(portType, net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

// Check /api/portchecker/check
// @Summary Check if host port is open for tcp or udp pconnections.
// @Description Check if host port is open for tcp or udp pconnections.
// @Tags Lookup
// @Accept json
// @Produce json,xml
// @Param "" body Request true "Request JSON"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Failure 429 {object} data.APIResponse
// @Router /api/portchecker/check [post]
func (API) Check(c *gin.Context) {
	response := &data.APIResponse{}
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.Message = "Invalid Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	if request.Port < 0 || request.Port > 65535 {
		response.Message = "Port should be between 0 and 65535"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	request.Host = strings.TrimSpace(request.Host)
	if strings.HasPrefix(request.Host, "192.168.") || request.Host == "127.0.0.1" ||
		request.Host == "0.0.0.0" || request.Host == "localhost" {
		response.Message = "Please don't do dirty things"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	request.PortType = strings.TrimSpace(strings.ToLower(request.PortType))

	if request.PortType == "tcp" || request.PortType == "udp" {
		redisKey := getRedisKey(c.ClientIP())
		if !db.RedisExist(redisKey) {
			db.RedisSet(redisKey, 1, 3*time.Second)
		} else {
			response.Message = "Too many requests today. Please take a break for 3 seconds."
			core.WriteResponse(c, http.StatusTooManyRequests, response)
			return
		}

		open := checkPort(request.Host, strconv.Itoa(request.Port), request.PortType)
		response.Data = open
		core.WriteResponse(c, http.StatusOK, response)
		return
	}
	response.Message = "Only tcp and udp are supported"
	core.WriteResponse(c, http.StatusBadRequest, response)
}

func getRedisKey(ip string) string {
	return fmt.Sprintf("APP_PORT_CHECKER_%s", ip)
}
