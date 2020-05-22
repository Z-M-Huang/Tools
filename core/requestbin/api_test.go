package requestbin

import (
	"os"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
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
		DatabaseConfig: &data.DatabaseConfiguration{
			ConnectionString: "./test.db",
			Driver:           "sqlite3",
		},
		RedisConfig: &data.RedisConfiguration{
			Addr: mr.Addr(),
		},
		Host: "localhost",
	}
	db.InitRedis()
}

func TestCreate(t *testing.T) {
	bin := create(false)

	assert.NotEmpty(t, bin)
	assert.NotEmpty(t, bin.ID)
	assert.NotEmpty(t, bin.URL)
	assert.Empty(t, bin.VerificationKey)

	data.Config.HTTPS = true
	privateBin := create(true)
	assert.NotEmpty(t, privateBin)
	assert.NotEmpty(t, privateBin.ID)
	assert.NotEmpty(t, privateBin.URL)
	assert.NotEmpty(t, privateBin.VerificationKey)
}

func TestGetRequestBinHistory(t *testing.T) {
	binData := create(false)

	redisBinData := GetRequestBinHistory(binData.ID)
	assert.NotEmpty(t, redisBinData)
}

func TestGetRequestBinHistoryFail(t *testing.T) {
	redisBinData := GetRequestBinHistory("123456")
	assert.Empty(t, redisBinData)
}
