package account

import (
	"os"
	"testing"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
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
		JwtKey:          []byte("CBYtDWTfRU5Pv7yULj46vm8ueZG7hbnq"),
		Host:            "localhost",
		ResourceVersion: "1",
		IsDebug:         true,
		HTTPS:           false,
		EnableSitemap:   true,
	}

	db.InitDB()
	db.InitRedis()
	InitGoogleOauth()
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

func TestIsTokenValid(t *testing.T) {
	tokenStr, expiresAt, err := generateJWTToken("Test", "test@example.com", "testUser", "https://localhost/imageURL")

	assert.Empty(t, err)
	assert.True(t, expiresAt.After(time.Now()))
	assert.NotEmpty(t, tokenStr)

	claim, err := isTokenValid(tokenStr)

	assert.Empty(t, err)
	assert.NotEmpty(t, claim)
	assert.NotEmpty(t, claim.Id)
}

func TestIsTokenValidFail(t *testing.T) {
	tokenStr := "123"

	claim, err := isTokenValid(tokenStr)

	assert.NotEmpty(t, err)
	assert.Empty(t, claim)
}
