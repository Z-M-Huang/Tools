package home

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/alicebob/miniredis"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	data.Config = &data.Configuration{
		JwtKey:          []byte("CBYtDWTfRU5Pv7yULj46vm8ueZG7hbnq"),
		Host:            "localhost",
		ResourceVersion: "1",
		IsDebug:         true,
		HTTPS:           false,
		EnableSitemap:   true,
	}
	data.DatabaseConfig = &data.DatabaseConfiguration{
		ConnectionString: "./test.db",
		Driver:           "sqlite3",
	}
	data.RedisConfig = &data.RedisConfiguration{
		Addr: mr.Addr(),
	}
	data.GoogleOauthConfig = &data.GoogleOauthConfiguration{
		ClientID:     "testClientID",
		ClientSecret: "testClientSecret",
	}

	db.InitDB()
	db.InitRedis()
}

func teardown() {
	err := db.Disconnect()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		err = os.Remove(data.DatabaseConfig.ConnectionString)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}
}

func loadTemplateRenderer() *template.Template {
	var err error
	t := template.New("")
	t.Funcs(template.FuncMap{"add": func(i, j int) int { return i + j }})
	t.Funcs(template.FuncMap{"mod": func(i, j int) int { return i % j }})
	t.Funcs(template.FuncMap{"nospace": func(i string) string {
		return strings.ReplaceAll(i, " ", "")
	}})
	t, err = t.ParseFiles(getAlltemplates("../../templates/")...)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	return t
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
