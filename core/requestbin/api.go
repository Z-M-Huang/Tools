package requestbin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//API request bin
type API struct{}

func create(private bool) *BinData {
	binData := &BinData{
		ID: strconv.FormatInt(time.Now().Unix(), 10),
	}
	if data.Config.HTTPS {
		binData.URL = "https://"
	} else {
		binData.URL = "http://"
	}

	binData.URL += data.Config.Host + "/api/request-bin/receive/" + binData.ID

	if private {
		binData.VerificationKey = utils.RandomString(30)
	}

	bytes, err := json.Marshal(binData)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	key := GetRequestBinKey(binData.ID)
	err = db.RedisSet(key, bytes, 24*time.Hour)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return binData
}

// CreateRequestBin /api/request-bin/Create
// @Summary Receive and visualize HTTP requests for any method.
// @Description Receive and visualize HTTP requests for any method.
// @Tags Web-Utils
// @Accept json
// @Produce json,xml
// @Param "" body CreateBinRequest true "Request JSON"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Router /api/request-bin/Create [post]
func (API) CreateRequestBin(c *gin.Context) {
	response := &data.APIResponse{}
	request := &CreateBinRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		core.WriteResponse(c, http.StatusBadRequest, &data.APIResponse{
			Message: "Invalid request.",
		})
		return
	}

	bin := create(request.IsPrivate)
	if bin == nil {
		core.WriteResponse(c, http.StatusInternalServerError, response)
		c.Abort()
		return
	}
	result := &CreateBinResponse{
		ID:              bin.ID,
		VerificationKey: bin.VerificationKey,
	}

	response.Data = result
	core.WriteResponse(c, 200, response)
}

// RequestIn /api/request-bin/receive/:id
// @Summary Take any request for request-bin
// @Description Take any request for request-bin
// @Tags Web-Utils
// @Accept json,xml
// @Produce json,xml
// @Param "" body object true "Any kind of request"
// @Param id path int true "id from create request"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Router /api/request-bin/receive/{id} [post]
func (API) RequestIn(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	binData := GetRequestBinHistory(id)
	if binData == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	history := &History{
		TimeReceived:        time.Now().UTC(),
		Method:              c.Request.Method,
		Proto:               c.Request.Proto,
		RemoteAddr:          c.Request.RemoteAddr,
		QueryStrings:        c.Request.URL.RawQuery,
		Headers:             make(map[string]string),
		Forms:               make(map[string]string),
		MultipartFormsFiles: make(map[string]string),
	}

	for key, val := range c.Request.Header {
		if strings.ToLower(key) != "cookie" {
			headerBytes, err := json.Marshal(val)
			if err != nil {
				utils.Logger.Error(err.Error())
				history.Headers[key] = "Unknown"
			} else {
				history.Headers[key] = string(headerBytes)
			}
		}
	}

	for _, cookie := range c.Request.Cookies() {
		history.Cookies = append(history.Cookies, cookie.String())
	}

	err := c.Request.ParseMultipartForm(200)
	if err == nil {
		for key, val := range c.Request.Form {
			formBytes, err := json.Marshal(val)
			if err != nil {
				utils.Logger.Error(err.Error())
				history.Forms[key] = "Unknown"
			} else {
				history.Forms[key] = string(formBytes)
			}
		}

		for key, val := range c.Request.MultipartForm.File {
			history.MultipartFormsFiles[key] = val[0].Filename
		}
	}

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	history.Body = string(bodyBytes)

	binData.History = append([]*History{history}, binData.History...)
	if len(binData.History) >= 20 {
		binData.History = binData.History[0:19]
	}
	key := GetRequestBinKey(binData.ID)
	err = db.RedisSetBytes(key, binData, 24*time.Hour)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(200, "Okay")
}

//GetRequestBinHistory get request history by id
func GetRequestBinHistory(id string) *BinData {
	key := GetRequestBinKey(id)
	data := &BinData{}
	err := db.RedisGet(key, data)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return data
}

//GetRequestBinKey request bin key
func GetRequestBinKey(id string) string {
	return fmt.Sprintf("APP_REQUEST_BIN_%s", id)
}
