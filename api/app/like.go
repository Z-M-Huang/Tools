package app

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//Like /app/like/:name
func Like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Only logged in user can access this
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)

	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	appCard := utils.GetApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	user := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: claim.Id,
	}).First(&user); db.RecordNotFound() || user == nil {
		utils.Logger.Sugar().Errorf("oh boy... There is a user doesn't found in database but have a token. Email: %s", claim.Id)
		response.Alert.IsDanger = true
		response.Alert.Message = "User not found"
		api.WriteResponse(w, response)
		return
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		response.Alert.IsDanger = true
		response.Alert.Message = "User not found"
		api.WriteResponse(w, response)
		return
	}

	found := false
	for _, likedApp := range user.LikedApps {
		if likedApp == appCard.Title {
			found = true
			break
		}
	}

	if !found {
		user.LikedApps = append(user.LikedApps, appCard.Title)
		if db := utils.DB.Save(user).Scan(&user); db.Error != nil {
			utils.Logger.Sugar().Errorf("failed to add liked app %s to user %s", appCard.Title, claim.Id)
		}
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
