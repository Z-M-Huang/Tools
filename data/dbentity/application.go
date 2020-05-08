package dbentity

import "github.com/jinzhu/gorm"

//Application list of applications
type Application struct {
	gorm.Model
	Name  string `gorm:"primary_key"`
	Usage uint64 `gorm:"not null;default:0"`
	Liked uint64 `gorm:"not null;default:0"`
}
