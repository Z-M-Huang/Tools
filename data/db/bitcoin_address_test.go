package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSaveBitcoinAddress(t *testing.T) {
	a := &User{
		Email:    "TestSaveBitcoinAddress@example.com",
		Username: "TestSaveBitcoinAddress",
	}
	assert.Empty(t, a.Save())
	address := "TestSaveBitcoinAddress"

	b := BitcoinAddress{
		Address: &address,
		UserID:  &a.ID,
	}
	assert.Empty(t, b.Save())
}

func TestSaveBitcoinAddressFail(t *testing.T) {
	a := &User{
		Email:    "TestSaveBitcoinAddressFail@example.com",
		Username: "TestSaveBitcoinAddressFail",
	}
	address := "TestSaveBitcoinAddressFail"
	assert.Empty(t, a.Save())
	b := []BitcoinAddress{
		{
			UserID: &a.ID,
		},
		{
			Address: &address,
		},
		{},
	}

	for _, i := range b {
		assert.NotEmpty(t, i.Save(), i.UserID)
	}
}

func TestSaveBitcoinAddressWithTx(t *testing.T) {
	a := &User{
		Email:    "TestSaveBitcoinAddressWithTx@example.com",
		Username: "TestSaveBitcoinAddressWithTx",
	}
	assert.Empty(t, a.Save())
	address := "TestSaveBitcoinAddressWithTx"

	b := BitcoinAddress{
		Address: &address,
		UserID:  &a.ID,
	}
	err := DoTransaction(func(tx *gorm.DB) error {
		return b.SaveWithTx(tx)
	})
	assert.Empty(t, err)
}

func TestSaveBitcoinAddressWithTxFail(t *testing.T) {
	a := &User{
		Email:    "TestSaveBitcoinAddressWithTxFail@example.com",
		Username: "TestSaveBitcoinAddressWithTxFail",
	}
	address := "TestSaveBitcoinAddressWithTxFail"
	assert.Empty(t, a.Save())
	b := []BitcoinAddress{
		{
			UserID: &a.ID,
		},
		{
			Address: &address,
		},
		{},
	}

	for _, i := range b {
		err := DoTransaction(func(tx *gorm.DB) error {
			return i.SaveWithTx(tx)
		})
		assert.NotEmpty(t, err)
	}
}

func TestFindBitcoinAddress(t *testing.T) {
	a := &User{
		Email:    "TestFindBitcoinAddress@example.com",
		Username: "TestFindBitcoinAddress",
	}
	assert.Empty(t, a.Save())
	address := "TestFindBitcoinAddress"
	b := BitcoinAddress{
		Address: &address,
		UserID:  &a.ID,
	}
	assert.Empty(t, b.Save())

	c := &BitcoinAddress{
		Address: &address,
	}
	assert.Empty(t, c.Find())
	assert.Equal(t, a.ID, *c.UserID)
}

func TestFindBitcoinAddressFail(t *testing.T) {
	address := "TestFindBitcoinAddressFail"
	c := &BitcoinAddress{
		Address: &address,
	}
	assert.NotEmpty(t, c.Find())
}

func TestFindBitcoinAddressWithTx(t *testing.T) {
	a := &User{
		Email:    "TestFindBitcoinAddressWithTx@example.com",
		Username: "TestFindBitcoinAddressWithTx",
	}
	assert.Empty(t, a.Save())
	address := "TestFindBitcoinAddressWithTx"
	b := BitcoinAddress{
		Address: &address,
		UserID:  &a.ID,
	}
	assert.Empty(t, b.Save())

	DoTransaction((func(tx *gorm.DB) error {
		c := &BitcoinAddress{
			Address: &address,
		}
		assert.Empty(t, c.FindWithTx(tx))
		assert.Equal(t, a.ID, *c.UserID)
		return nil
	}))
}

func TestFindBitcoinAddressWithTxFail(t *testing.T) {
	address := "TestFindBitcoinAddressWithTxFail"
	DoTransaction((func(tx *gorm.DB) error {
		c := &BitcoinAddress{
			Address: &address,
		}
		assert.NotEmpty(t, c.FindWithTx(tx))
		return nil
	}))
}
