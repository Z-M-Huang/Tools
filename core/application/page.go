package application

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/core/requestbin"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//Page logic
type Page struct{}

//RenderApplicationPage renders /app/:name
func (Page) RenderApplicationPage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.PageResponse)

	name := c.Param("name")

	if name == "" {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	appCard := GetApplicationsByName(name)
	if appCard == nil {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}
	response.Header.Title = appCard.Title + " - Fun Apps"
	response.Header.Description = appCard.Description

	AddApplicationUsage(appCard)
	response.Data = loadAppSpecificData(c, appCard.Name)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, appCard.TemplateName, response)
}

func loadAppSpecificData(c *gin.Context, appName string) interface{} {
	switch appName {
	case "request-bin":
		return requestbin.LoadRequestBinData(c)
	case "port-checker":
		return c.ClientIP()
	}
	return nil
}
