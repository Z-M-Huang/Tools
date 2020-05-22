package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	a := &User{
		Email:    "testSave@example.com",
		Username: "testSave",
	}
	err := a.Save()

	assert.Empty(t, err)
}

func TestSaveUserFail(t *testing.T) {
	a := &User{
		Email:    "testSaveFail@example.com",
		Username: "testSaveFail",
	}
	assert.Empty(t, a.Save())
	b := &User{
		Email:    a.Email,
		Username: "testSaveFailb",
	}
	err := b.Save()
	//Duplicate email
	assert.NotEmpty(t, err)

	c := &User{
		Email:    "testSaveFailc@example.com",
		Username: a.Username,
	}
	err = c.Save()
	//Duplicate Username
	assert.NotEmpty(t, err)
}

func TestSaveUserWithTx(t *testing.T) {
	a := &User{
		Email:    "testSaveWithTx@example.com",
		Username: "testSaveWithTx",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.SaveWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestSaveUserWithTxFail(t *testing.T) {
	a := &User{
		Email:    "testSaveWithTxFail@example.com",
		Username: "testSaveWithTxFail",
	}
	assert.Empty(t, a.Save())
	b := &User{
		Email:    a.Email,
		Username: "testSaveWithTxFailb",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return b.SaveWithTx(tx)
	})
	assert.NotEmpty(t, err)

	c := &User{
		Email:    "testSaveWithTxFailc@example.com",
		Username: a.Username,
	}
	//Duplicate Username
	err = DoTransaction(func(tx *gorm.DB) error {
		return c.SaveWithTx(tx)
	})

	assert.NotEmpty(t, err)
}

func TestFindUser(t *testing.T) {
	a := &User{
		Email:    "testFind@example.com",
		Username: "testFind",
	}
	assert.Empty(t, a.Save())

	err := a.Find()

	assert.Empty(t, err)
}

func TestFindUserFail(t *testing.T) {
	a := &User{
		Email: "testFindFail@example.com",
	}

	err := a.Find()

	assert.NotEmpty(t, err)
}

func TestFindUserWithTx(t *testing.T) {
	a := &User{
		Email:    "testFindWithTx@example.com",
		Username: "testFindWithTx",
	}
	assert.Empty(t, a.Save())

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestFindUserWithTxFail(t *testing.T) {
	a := &User{
		Email: "testFindWithTxFail@example.com",
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.NotEmpty(t, err)
}
