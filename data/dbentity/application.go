package dbentity

import "github.com/jinzhu/gorm"

//Application list of applications
type Application struct {
	gorm.Model
	Name  string `gorm:"primary_key"`
	Usage int64  `gorm:"not null;default:0"`
	Liked int64  `gorm:"not null;default:0"`
}
