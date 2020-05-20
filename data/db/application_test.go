package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSaveApp(t *testing.T) {
	a := &Application{
		Name: "testApp",
	}
	err := a.Save()

	assert.Empty(t, err)
}

func TestSaveAppFail(t *testing.T) {
	a := &Application{
		Name: "testAppFail",
	}
	assert.Empty(t, a.Save())
	b := &Application{
		Name: a.Name,
	}
	err := b.Save()
	//Duplicate
	assert.NotEmpty(t, err)
}

func TestSaveAppWithTx(t *testing.T) {
	a := &Application{
		Name: "testAppSaveWithTx",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.SaveWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestSaveAppWithTxFail(t *testing.T) {
	a := &Application{
		Name: "testAppSaveWithTxFail",
	}
	assert.Empty(t, a.Save())
	b := &Application{
		Name: a.Name,
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return b.SaveWithTx(tx)
	})

	assert.NotEmpty(t, err)
}

func TestFindApp(t *testing.T) {
	a := &Application{
		Name: "testFind",
	}
	assert.Empty(t, a.Save())

	err := a.Find()

	assert.Empty(t, err)
}

func TestFindAppFail(t *testing.T) {
	a := &Application{
		Name: "testFindFail",
	}

	err := a.Find()

	assert.NotEmpty(t, err)
}

func TestFindAppWithTx(t *testing.T) {
	a := &Application{
		Name: "testAppFindWithTx",
	}
	assert.Empty(t, a.Save())

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestFindAppWithTxFail(t *testing.T) {
	a := &Application{
		Name: "testAppFindWithTxFail",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.NotEmpty(t, err)
}
