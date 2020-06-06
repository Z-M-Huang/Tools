package requestbin

import (
	"os"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/alicebob/miniredis"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	os.Exit(ret)
}

func setup() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	data.Config = &data.Configuration{
		Host: "localhost",
	}
	data.DatabaseConfig = &data.DatabaseConfiguration{
		ConnectionString: "./test.db",
		Driver:           "sqlite3",
	}
	data.RedisConfig = &data.RedisConfiguration{
		Addr: mr.Addr(),
	}
	db.InitRedis()
}
