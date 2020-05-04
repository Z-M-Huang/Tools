package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

var tplt *template.Template

func init() {
	var err error
	tplt = template.New("")
	getTemplateFuncs()
	tplt, err = tplt.ParseFiles(getAlltemplates("templates/")...)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageData := &webdata.PageData{}
	var cardList []*webdata.AppCardList
	for i := 0; i < 5; i++ {
		cardCategory := &webdata.AppCardList{
			Category: strconv.Itoa(i),
		}
		for j := 0; j < 10; j++ {
			cardCategory.AppCards = append(cardCategory.AppCards, &webdata.AppCard{
				Link:        "#",
				Title:       fmt.Sprintf("Card Title %d-%d", i, j),
				Description: "Some quick example text to build on the card title and make up the bulk of the card's content.",
				Up:          10000,
				Saved:       10000,
			})
		}
		cardList = append(cardList, cardCategory)
	}
	pageData.ContentData = cardList
	claim, err := api.GetClaimFromToken(w, r)
	if err == nil && claim != nil {
		pageData.LoginInfo = webdata.LoginData{
			Username: claim.Subject,
			ImageURL: claim.ImageURL,
		}
	}
	tplt.ExecuteTemplate(w, "homepage.gohtml", pageData)
}

func loginPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tplt.ExecuteTemplate(w, "login.gohtml", &webdata.PageData{})
}

func accountPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageData := &webdata.PageData{}
	claim, err := api.GetClaimFromToken(w, r)
	if err != nil {
		pageData.AlertInfo.IsDanger = true
		pageData.AlertInfo.Message = err.Error()
	} else {
		pageData.LoginInfo = webdata.LoginData{
			Username: claim.Subject,
			ImageURL: claim.ImageURL,
		}
		user, err := api.GetUserInfoFromDB(claim.Id)
		if err != nil {
			pageData.AlertInfo.IsDanger = true
			pageData.AlertInfo.Message = err.Error()
		} else {
			pageData.ContentData = user
		}
	}
	tplt.ExecuteTemplate(w, "account.gohtml", pageData)
}

func main() {

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/vendor/*filepath", http.Dir("node_modules/"))

	router.GET("/", homePage)
	router.GET("/login", loginPage)
	router.POST("/login", api.Login)
	router.GET("/google_login", api.GoogleLogin)
	router.GET("/google_oauth", api.GoogleCallback)

	utils.Logger.Fatal(http.ListenAndServe(":80", router).Error())
}

func getAlltemplates(inputPath string) []string {
	var ret []string
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if path != inputPath && info.IsDir() {
			ret = append(ret, getAlltemplates(path)...)
		} else if strings.Contains(info.Name(), ".gohtml") {
			ret = append(ret, path)
		}
		return nil
	})
	return ret
}

func getTemplateFuncs() {
	tplt.Funcs(template.FuncMap{"add": func(i, j int) int { return i + j }})
	tplt.Funcs(template.FuncMap{"mod": func(i, j int) int { return i % j }})
}
