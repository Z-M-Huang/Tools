package app

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//EncodeDecode /api/string/encodedecode
func EncodeDecode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	request := &application.StringEncodeDecodeRequest{}
	var result []string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid lookup request.",
		})
		response.Data = []string{response.Header.Alert.Message}
		api.WriteResponse(w, response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid lookup request.",
		})
		api.WriteResponse(w, response)
		return
	}

	request.Action = strings.TrimSpace(request.Action)
	request.Type = strings.TrimSpace(request.Type)
	lines := strings.Split(request.RequestString, "\r\n")

	if request.Action == "" || (request.Action != "encode" && request.Action != "decode") {
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Invalid action code.",
		})
		response.Data = []string{response.Header.Alert.Message}
		api.WriteResponse(w, response)
		return
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
					response.SetAlert(&data.AlertData{
						IsWarning: true,
						Message:   fmt.Sprintf("Cannot decode string requested %s", err.Error()),
					})
					response.Data = []string{response.Header.Alert.Message}
					api.WriteResponse(w, response)
					return
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
					response.SetAlert(&data.AlertData{
						IsWarning: true,
						Message:   fmt.Sprintf("Cannot decode string requested %s", err.Error()),
					})
					response.Data = []string{response.Header.Alert.Message}
					api.WriteResponse(w, response)
					return
				}
				result = append(result, string(unescaped))
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
					response.SetAlert(&data.AlertData{
						IsDanger: true,
						Message:  fmt.Sprintf("Cannot decode string requested %s", err.Error()),
					})
					response.Data = []string{response.Header.Alert.Message}
					api.WriteResponse(w, response)
					return
				}
				result = append(result, unescaped)
			}
		}
	default:
		response.SetAlert(&data.AlertData{
			IsWarning: true,
			Message:   "Invalid type request.",
		})
		response.Data = []string{response.Header.Alert.Message}
		api.WriteResponse(w, response)
		return
	}

	response.Data = result
	api.WriteResponse(w, response)
}
