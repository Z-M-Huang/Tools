package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//HomePage home page /
func HomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	response.Data = utils.AppList
	response.Header.Title = "Fun Apps"
	utils.Templates.ExecuteTemplate(w, "homepage.gohtml", response)
}
