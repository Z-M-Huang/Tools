package db

import "github.com/jinzhu/gorm"

//Application list of applications
type Application struct {
	gorm.Model
	Name  string `gorm:"unique;unique_index"`
	Usage uint64 `gorm:"not null;default:0"`
	Liked uint64 `gorm:"not null;default:0"`
}

//Find populate current object
func (a *Application) Find() error {
	if db := dbContext.Where(*a).First(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//FindWithTx populate current object with transaction
func (a *Application) FindWithTx(tx *gorm.DB) error {
	if db := tx.Where(*a).First(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save application
func (a *Application) Save() error {
	if db := dbContext.Save(a).Scan(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//SaveWithTx save application with transaction
func (a *Application) SaveWithTx(tx *gorm.DB) error {
	if db := tx.Save(a).Scan(&a); db.Error != nil {
		return db.Error
	}
	return nil
}
