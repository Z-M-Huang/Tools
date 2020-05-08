package pages

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//SignupPage /signup
func SignupPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim).IsNil() {
		utils.Templates.ExecuteTemplate(w, "signup.gohtml", data.Response{})
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

//LoginPage /login
func LoginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	if r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim).IsNil() {
		utils.Templates.ExecuteTemplate(w, "login.gohtml", response)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

//AccountPage /account requires claim
func AccountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	response.Login = data.LoginData{
		Username: claim.Subject,
		ImageURL: claim.ImageURL,
	}
	user, err := GetUserInfoFromDB(claim.Id)
	if err != nil {
		response.Alert.IsDanger = true
		response.Alert.Message = err.Error()
	} else {
		response.Data = webdata.AccountPageData{
			HasPassword: user.Password != "",
		}
	}
	utils.Templates.ExecuteTemplate(w, "account.gohtml", response)
}

//GetUserInfoFromDB get user info from database
func GetUserInfoFromDB(emailAddress string) (*dbentity.User, error) {
	user := &dbentity.User{}
	if db := utils.DB.Where(dbentity.User{
		Email: emailAddress,
	}).First(&user); db.RecordNotFound() {
		return nil, fmt.Errorf("User not found associated with email: %s", emailAddress)
	} else if db.Error != nil {
		utils.Logger.Error(db.Error.Error())
		return nil, errors.New("Internal Error: failed to get user info, please try again later")
	}
	return user, nil
}
