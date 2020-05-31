package db

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//User database entity
type User struct {
	gorm.Model
	Username   string `gorm:"unique_index"`
	Password   string
	Email      string `gorm:"primary_key;unique"`
	GoogleID   string
	LikedApps  pq.StringArray `gorm:"type:text[]"`
	ShortLinks []*ShortLink
}

//Find populate current object
func (u *User) Find() error {
	if db := dbContext.Where(*u).First(&u); db.Error != nil {
		return db.Error
	}
	return nil
}

//FindWithTx populate current object with transaction
func (u *User) FindWithTx(tx *gorm.DB) error {
	if db := tx.Where(*u).First(&u); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save current user
func (u *User) Save() error {
	if db := dbContext.Save(u).Scan(&u); db.Error != nil {
		return db.Error
	}
	return nil
}

//SaveWithTx save current user with transaction
func (u *User) SaveWithTx(tx *gorm.DB) error {
	if db := tx.Save(u).Scan(&u); db.Error != nil {
		return db.Error
	}
	return nil
}
