package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var tplt *template.Template
var logger *zap.Logger

func init() {
	initLogger()
	var err error
	tplt, err = template.ParseFiles(getAlltemplates("templates/")...)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tplt.ExecuteTemplate(w, "homepage.gohtml", nil)
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
