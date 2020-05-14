package pages

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/logic"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//RenderApplicationPage renders /app/:name
func RenderApplicationPage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)

	name := c.Param("name")

	if name == "" {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	appCard := applicationlogic.GetApplicationsByName(name)
	if appCard == nil {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}
	response.Header.Title = appCard.Title + " - Fun Apps"
	response.Header.Description = appCard.Description

	usedApps, err := applicationlogic.GetApplicationUsed(c.Request)
	if err == nil {
		exists := false
		for _, str := range usedApps {
			if str == appCard.Title {
				exists = true
				break
			}
		}
		if !exists {
			addApplicationUsage(appCard)
			usedApps = append(usedApps, appCard.Title)
			usedAppsBytes, err := json.Marshal(usedApps)
			encoded := base64.StdEncoding.EncodeToString(usedAppsBytes)
			if err != nil {
				utils.Logger.Error(err.Error())
			} else {
				logic.SetCookie(c, utils.UsedTokenKey, string(encoded), time.Date(2199, time.December, 31, 23, 59, 59, 0, time.UTC))
			}
		}
	} else {
		utils.Logger.Error(err.Error())
	}

	utils.Templates.ExecuteTemplate(c.Writer, appCard.TemplateName, response)
}

func addApplicationUsage(app *webdata.AppCard) {
	app.AmountUsed++
	dbApp := &dbentity.Application{
		Name: app.Title,
	}
	err := applicationlogic.Find(utils.DB, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	dbApp.Usage++
	err = applicationlogic.Save(utils.DB, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}
