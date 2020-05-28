package emailmmssms

import (
	"os"
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
		DatabaseConfig: &data.DatabaseConfiguration{
			ConnectionString: "./test.db",
			Driver:           "sqlite3",
		},
		RedisConfig: &data.RedisConfiguration{
			Addr: mr.Addr(),
		},
		GoogleOauthConfig: &data.GoogleOauthConfiguration{
			ClientID:     "testClientID",
			ClientSecret: "testClientSecret",
		},
		EmailConfig:     &data.EmailConfiguration{},
		JwtKey:          []byte("CBYtDWTfRU5Pv7yULj46vm8ueZG7hbnq"),
		Host:            "localhost",
		ResourceVersion: "1",
		IsDebug:         true,
		HTTPS:           false,
		EnableSitemap:   true,
	}

	db.InitDB()
	db.InitRedis()
}

func teardown() {
	err := db.Disconnect()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		err = os.Remove(data.Config.DatabaseConfig.ConnectionString)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}
}
