package db

import (
	"os"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/alicebob/miniredis"
	"github.com/jinzhu/gorm"
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
			ConnectionString: "./testapplication.db",
			Driver:           "sqlite3",
		},
		RedisConfig: &data.RedisConfiguration{
			Addr: mr.Addr(),
		},
	}

	InitDB()
	InitRedis()
}

func teardown() {
	err := Disconnect()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		err = os.Remove(data.Config.DatabaseConfig.ConnectionString)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}
}

func TestSave(t *testing.T) {
	a := &Application{
		Name: "testApp",
	}
	err := a.Save()

	assert.Empty(t, err)
}

func TestSaveFail(t *testing.T) {
	a := &Application{
		Name: "testAppFail",
	}
	a.Save()
	b := &Application{
		Name: a.Name,
	}
	err := b.Save()
	//Duplicate
	assert.NotEmpty(t, err)
}

func TestSaveWithTx(t *testing.T) {
	a := &Application{
		Name: "testAppSaveWithTx",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.SaveWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestSaveWithTxFail(t *testing.T) {
	a := &Application{
		Name: "testAppSaveWithTxFail",
	}
	a.Save()
	b := &Application{
		Name: a.Name,
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return b.SaveWithTx(tx)
	})

	assert.NotEmpty(t, err)
}

func TestFind(t *testing.T) {
	a := &Application{
		Name: "testFind",
	}
	a.Save()

	err := a.Find()

	assert.Empty(t, err)
}

func TestFindFail(t *testing.T) {
	a := &Application{
		Name: "testFindFail",
	}

	err := a.Find()

	assert.NotEmpty(t, err)
}

func TestFindWithTx(t *testing.T) {
	a := &Application{
		Name: "testAppFindWithTx",
	}
	a.Save()

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestFindWithTxFail(t *testing.T) {
	a := &Application{
		Name: "testAppFindWithTxFail",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.NotEmpty(t, err)
}
