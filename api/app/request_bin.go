package app

import (
	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/apidata/application"
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

	bin := applicationlogic.NewRequestBinHistory(c, request.IsPrivate)
	if bin == nil {
		api.WriteUnexpectedError(c, response)
		c.Abort()
		return
	}
	result := &application.CreateBinResponse{
		URL:             bin.URL,
		VerificationKey: bin.VerificationKey,
	}

	response.Data = result
	api.WriteResponse(c, 200, response)
}
