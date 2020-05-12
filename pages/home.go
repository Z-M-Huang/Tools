package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//HomePage home page /
func HomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	if !(claim == nil) {
		user := &dbentity.User{
			Email: claim.Id,
		}
		err := userlogic.Find(utils.DB, user)
		if err == nil {
			if len(user.LikedApps) > 0 {
				response.Data = applicationlogic.GetApplicationWithLiked(user)
			} else {
				response.Data = utils.AppList
			}
		} else {
			response.Data = utils.AppList
		}
	} else {
		response.Data = utils.AppList
	}

	response.Header.Title = "Fun Apps"
	utils.Templates.ExecuteTemplate(w, "homepage.gohtml", response)
}
