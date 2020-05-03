package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var tplt *template.Template
var logger *zap.Logger

func init() {
	tplt = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tplt.ExecuteTemplate(w, "homepage.gohtml", nil)
}

func main() {
	initLogger()

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
