package db

import (
	"github.com/jinzhu/gorm"
)

//ShortLink database entity
type ShortLink struct {
	gorm.Model
	Link  *string `gorm:"not null"`
	Usage uint64  `gorm:"not null;default:0"`
	User  *User
}

//Find populate current object
func (s *ShortLink) Find() error {
	if db := dbContext.Where(*s).First(&s); db.Error != nil {
		return db.Error
	}
	return nil
}

//FindWithTx populate current object with transaction
func (s *ShortLink) FindWithTx(tx *gorm.DB) error {
	if db := tx.Where(*s).First(&s); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save current object
func (s *ShortLink) Save() error {
	if db := dbContext.Save(s).Scan(&s); db.Error != nil {
		return db.Error
	}
	return nil
}

//SaveWithTx save current object with transaction
func (s *ShortLink) SaveWithTx(tx *gorm.DB) error {
	if db := tx.Save(s).Scan(&s); db.Error != nil {
		return db.Error
	}
	return nil
}
