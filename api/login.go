package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Login request
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.Redirect(w, r, "/", 301)
}
