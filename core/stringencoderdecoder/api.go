package stringencoderdecoder

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//API string encoder decoder
type API struct{}

func encodeDecode(request *Request) (int, *data.APIResponse) {
	response := &data.APIResponse{}
	var result []string
	request.Action = strings.TrimSpace(request.Action)
	request.Type = strings.TrimSpace(request.Type)
	lines := strings.Split(request.RequestString, "\r\n")

	if request.Action == "" || (request.Action != "encode" && request.Action != "decode") {
		response.Message = "Invalid action code."
		response.Data = "Invalid action code."
		return http.StatusBadRequest, response
	}
	switch request.Type {
	case "Base32":
		if request.Action == "encode" {
			for _, line := range lines {
				result = append(result, base32.StdEncoding.EncodeToString([]byte(line)))
			}
		} else {
			for _, line := range lines {
				unescaped, err := base32.StdEncoding.DecodeString(line)
				if err != nil {
					response.Message = fmt.Sprintf("Cannot decode string requested %s", err.Error())
					response.Data = response.Message
					return http.StatusBadRequest, response
				}
				result = append(result, string(unescaped))
			}
		}
	case "Base64":
		if request.Action == "encode" {
			for _, line := range lines {
				result = append(result, base64.StdEncoding.EncodeToString([]byte(line)))
			}
		} else {
			for _, line := range lines {
				unescaped, err := base64.StdEncoding.DecodeString(line)
				if err != nil {
					response.Message = fmt.Sprintf("Cannot decode string requested %s", err.Error())
					response.Data = response.Message
					return http.StatusBadRequest, response
				}
				result = append(result, string(unescaped))
			}
		}
	case "Binary":
		if request.Action == "encode" {
			for _, line := range lines {
				temp := ""
				for _, c := range line {
					bArr := make([]byte, utf8.RuneLen(c))
					utf8.EncodeRune(bArr, c)
					for _, b := range bArr {
						temp += fmt.Sprintf("%.8b ", b)
					}
				}
				result = append(result, temp)
			}
		} else {
			for _, line := range lines {
				b := make([]byte, 0)
				for _, s := range strings.Fields(line) {
					n, _ := strconv.ParseUint(s, 2, len(s))
					b = append(b, byte(n))
				}
				result = append(result, string(b))
			}
		}
	case "URL":
		if request.Action == "encode" {
			for _, line := range lines {
				result = append(result, url.QueryEscape(line))
			}
		} else {
			for _, line := range lines {
				unescaped, err := url.QueryUnescape(line)
				if err != nil {
					response.Message = fmt.Sprintf("Cannot decode string requested %s", err.Error())
					response.Data = response.Message
					return http.StatusBadRequest, response
				}
				result = append(result, unescaped)
			}
		}
	default:
		response.Message = "Invalid type request."
		response.Data = response.Message
		return http.StatusBadRequest, response
	}
	response.Data = result
	return http.StatusOK, response
}

//EncodeDecode /api/string/encodedecode
func (API) EncodeDecode(c *gin.Context) {
	response := &data.APIResponse{}
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Invalid request."
		response.Data = "Invalid request."
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	status, response := encodeDecode(request)

	core.WriteResponse(c, status, response)
}
