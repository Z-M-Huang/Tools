package db

import (
	"os"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
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

	data.Config = &data.Configuration{}
	data.DatabaseConfig = &data.DatabaseConfiguration{
		ConnectionString: "./test.db",
		Driver:           "sqlite3",
	}
	data.RedisConfig = &data.RedisConfiguration{
		Addr: mr.Addr(),
	}
	InitDB()
	InitRedis()
}

func teardown() {
	err := Disconnect()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		err = os.Remove(data.DatabaseConfig.ConnectionString)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}
}
