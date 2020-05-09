package app

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//Like /app/like/:name
func Like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	appCard := utils.GetApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	addApplicationLike(appCard)
}

func addApplicationLike(app *webdata.AppCard) {
	app.AmountLiked++
	dbApp := &dbentity.Application{}
	if db := utils.DB.Where(dbentity.Application{
		Name: app.Title,
	}).First(&dbApp); db.Error != nil {
		utils.Logger.Error(db.Error.Error())
	} else if dbApp != nil {
		dbApp.Liked++
		if savedb := utils.DB.Save(dbApp).Scan(&dbApp); savedb.Error != nil {
			utils.Logger.Error(savedb.Error.Error())
		}
	} else {
		utils.Logger.Sugar().Errorf("failed to update application usage for %s", app.Title)
	}
}
