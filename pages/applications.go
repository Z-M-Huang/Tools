package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
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

	appCard := getApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	addApplicationUsage(appCard)
	utils.Templates.ExecuteTemplate(w, appCard.TemplateName, response)
}

func getApplicationsByName(name string) *webdata.AppCard {
	for _, category := range utils.AppList {
		for _, app := range category.AppCards {
			if app.Name == name {
				return app
			}
		}
	}
	return nil
}

func addApplicationUsage(app *webdata.AppCard) {
	app.Usage++
	dbApp := &dbentity.Application{}
	if db := utils.DB.Where(dbentity.Application{
		Name: app.Title,
	}).First(&dbApp); db.Error != nil {
		utils.Logger.Error(db.Error.Error())
	} else if dbApp != nil {
		dbApp.Usage++
		if savedb := utils.DB.Save(dbApp).Scan(&dbApp); savedb.Error != nil {
			utils.Logger.Error(savedb.Error.Error())
		}
	} else {
		utils.Logger.Sugar().Errorf("failed to update application usage for %s", app.Title)
	}
}
