package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var tplt *template.Template
var logger *zap.Logger

func init() {
	initLogger()
	var err error
	tplt = template.New("")
	getTemplateFuncs()
	tplt, err = tplt.ParseFiles(getAlltemplates("templates/")...)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cardList []*data.AppCardList
	for i := 0; i < 5; i++ {
		cardCategory := &data.AppCardList{
			Category: strconv.Itoa(i),
		}
		for j := 0; j < 10; j++ {
			cardCategory.AppCards = append(cardCategory.AppCards, &data.AppCard{
				Link:        "#",
				Title:       fmt.Sprintf("Card Title %d-%d", i, j),
				Description: "Some quick example text to build on the card title and make up the bulk of the card's content.",
				Up:          10000,
				Saved:       10000,
			})
		}
		cardList = append(cardList, cardCategory)
	}

	tplt.ExecuteTemplate(w, "homepage.gohtml", cardList)
}

func main() {

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/vendor/*filepath", http.Dir("node_modules/"))

	router.GET("/", homePage)

	logger.Fatal(http.ListenAndServe(":80", router).Error())
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	logger, _ = config.Build()
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
