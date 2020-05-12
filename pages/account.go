package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//SignupPage /signup
func SignupPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim).IsNil() {
		response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
		response.Header.Title = "Signup - Fun Apps"
		utils.Templates.ExecuteTemplate(w, "signup.gohtml", response)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

//LoginPage /login
func LoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	response.Header.Title = "Login - Fun Apps"
	redirectURL, ok := r.URL.Query()["redirect"]
	if ok && len(redirectURL) > 0 {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Please login first.",
		})
	}
	if r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim).IsNil() {
		utils.Templates.ExecuteTemplate(w, "login.gohtml", response)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

//AccountPage /account requires claim
func AccountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	response.Header.Title = "Account - Fun Apps"
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	user := &dbentity.User{
		Email: claim.Id,
	}
	err := userlogic.Find(utils.DB, user)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
		})
	} else {
		response.Data = webdata.AccountPageData{
			HasPassword: user.Password != "",
		}
	}
	utils.Templates.ExecuteTemplate(w, "account.gohtml", response)
}
