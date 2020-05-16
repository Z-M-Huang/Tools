package utils

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/jinzhu/gorm/dialects/mssql" //supporting packages
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var onceRedis sync.Once

//Templates page templates
var Templates *template.Template

//Logger global Logger
var Logger *zap.Logger

func init() {
	onceRedis.Do(func() {
		initLogger()
		initTemplates()
	})
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	Logger, _ = config.Build()
}

func initTemplates() {
	var err error
	Templates = template.New("")
	getTemplateFuncs()
	Templates, err = Templates.ParseFiles(getAlltemplates("templates/")...)
	if err != nil {
		Logger.Fatal(err.Error())
	}
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
	Templates.Funcs(template.FuncMap{"add": func(i, j int) int { return i + j }})
	Templates.Funcs(template.FuncMap{"mod": func(i, j int) int { return i % j }})
	Templates.Funcs(template.FuncMap{"nospace": func(i string) string {
		return strings.ReplaceAll(i, " ", "")
	}})
}
