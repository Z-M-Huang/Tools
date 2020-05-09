package applicationlogic

import (
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/jinzhu/gorm"
)

//Find find db application
func Find(tx *gorm.DB, a *dbentity.Application) error {
	if db := tx.Where(*a).First(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save application
func Save(tx *gorm.DB, a *dbentity.Application) error {
	if db := tx.Save(a).Scan(&a); db.Error != nil {
		return db.Error
	}
	return nil
}
