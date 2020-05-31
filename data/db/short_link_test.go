package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSaveShortLink(t *testing.T) {
	link := "https://example.com/123"
	a := &ShortLink{
		Link: &link,
	}
	err := a.Save()

	assert.Empty(t, err)
}

func TestSaveShortLinkFail(t *testing.T) {
	a := &ShortLink{}
	assert.NotEmpty(t, a.Save())
}

func TestSaveShortLinkWithTx(t *testing.T) {
	link := "https://example.com/123"
	a := &ShortLink{
		Link: &link,
	}
	err := DoTransaction(func(tx *gorm.DB) error {
		return a.SaveWithTx(tx)
	})
	assert.Empty(t, err)
}

func TestSaveShortLinkWithTxFail(t *testing.T) {
	a := &ShortLink{}
	err := DoTransaction(func(tx *gorm.DB) error {
		return a.SaveWithTx(tx)
	})
	assert.NotEmpty(t, err)
}

func TestFindShortLink(t *testing.T) {
	link := "https://example.com/123"
	a := &ShortLink{
		Link: &link,
	}
	assert.Empty(t, a.Save())

	err := a.Find()

	assert.Empty(t, err)
}

func TestFindShortLinkFail(t *testing.T) {
	a := &ShortLink{
		Model: gorm.Model{
			ID: 999999,
		},
	}

	err := a.Find()

	assert.NotEmpty(t, err)
}

func TestFindShortLinkWithTx(t *testing.T) {
	link := "https://example.com/123"
	a := &ShortLink{
		Link: &link,
	}
	assert.Empty(t, a.Save())

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.Empty(t, err)
}

func TestFindShortLinkWithTxFail(t *testing.T) {
	a := &ShortLink{
		Model: gorm.Model{
			ID: 999999,
		},
	}

	err := DoTransaction(func(tx *gorm.DB) error {
		return a.FindWithTx(tx)
	})

	assert.NotEmpty(t, err)
}
