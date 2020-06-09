package db

import "github.com/jinzhu/gorm"

//BitcoinAddress bitcoin address entity
type BitcoinAddress struct {
	gorm.Model
	Address *string `gorm:"not null;unique"`
	UserID  *uint   `gorm:"not null"`
	User    User    `gorm:"foreightkey:UserID"`
}

//Find populate current object
func (b *BitcoinAddress) Find() error {
	if db := dbContext.Where(*b).First(&b); db.Error != nil {
		return db.Error
	}
	return nil
}

//FindWithTx populate current object with transaction
func (b *BitcoinAddress) FindWithTx(tx *gorm.DB) error {
	if db := tx.Where(*b).First(&b); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save current user
func (b *BitcoinAddress) Save() error {
	if db := dbContext.Save(b).Scan(&b); db.Error != nil {
		return db.Error
	}
	return nil
}

//SaveWithTx save current user with transaction
func (b *BitcoinAddress) SaveWithTx(tx *gorm.DB) error {
	if db := tx.Save(b).Scan(&b); db.Error != nil {
		return db.Error
	}
	return nil
}
