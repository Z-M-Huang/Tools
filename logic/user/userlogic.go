package userlogic

import (
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/jinzhu/gorm"
)

//Find find db user
func Find(tx *gorm.DB, u *dbentity.User) error {
	if db := tx.Where(*u).First(&u); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save current user
func Save(tx *gorm.DB, u *dbentity.User) error {
	if db := tx.Save(u).Scan(&u); db.Error != nil {
		return db.Error
	}
	return nil
}
