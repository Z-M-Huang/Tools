package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
	webData "github.com/Z-M-Huang/Tools/data/webdata/application"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//CreateRequestBin /api/request-bin/Create
func CreateRequestBin(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)
	request := &application.CreateBinRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Invalid request.",
		})
		api.WriteResponse(c, 200, response)
		return
	}

	bin := applicationlogic.NewRequestBinHistory(request.IsPrivate)
	if bin == nil {
		api.WriteUnexpectedError(c, response)
		c.Abort()
		return
	}
	result := &application.CreateBinResponse{
		ID:              bin.ID,
		VerificationKey: bin.VerificationKey,
	}

	response.Data = result
	api.WriteResponse(c, 200, response)
}

//RequestIn /api/request-bin/receive/:id
func RequestIn(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data := applicationlogic.GetRequestBinHistory(id)
	if data == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	history := &webData.RequestHistory{
		TimeReceived: time.Now().UTC(),
		Method:       c.Request.Method,
		Proto:        c.Request.Proto,
		RemoteAddr:   c.Request.RemoteAddr,
		QueryStrings: c.Request.URL.RawQuery,
	}

	headerBytes, err := json.Marshal(c.Request.Header)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	history.Headers = string(headerBytes)

	formBytes, err := json.Marshal(c.Request.Form)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	history.Forms = string(formBytes)

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	history.Body = string(bodyBytes)

	data.History = append([]*webData.RequestHistory{history}, data.History...)
	if len(data.History) >= 20 {
		data.History = data.History[0:19]
	}
	key := applicationlogic.GetRequestBinKey(data.ID)
	bytes, err := json.Marshal(data)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = utils.RedisClient.Set(key, bytes, 24*time.Hour).Err()
	if err != nil {
		utils.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(200, "Okay")
}
