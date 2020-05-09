package pages

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//RenderApplicationPage renders /app/:name
func RenderApplicationPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)

	name := ps.ByName("name")

	if name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	appCard := applicationlogic.GetApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
	response.Header.Title = appCard.Title + " - Fun Apps"

	usedApps, err := applicationlogic.GetApplicationUsed(r)
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
				utils.SetCookie(w, utils.UsedTokenKey, string(encoded), time.Date(2199, time.December, 31, 23, 59, 59, 0, time.UTC))
			}
		}
	} else {
		utils.Logger.Error(err.Error())
	}

	utils.Templates.ExecuteTemplate(w, appCard.TemplateName, response)
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
