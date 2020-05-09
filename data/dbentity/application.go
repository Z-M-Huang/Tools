package dbentity

import "github.com/jinzhu/gorm"

//Application list of applications
type Application struct {
	gorm.Model
	Name  string `gorm:"unique;unique_index"`
	Usage uint64 `gorm:"not null;default:0"`
	Liked uint64 `gorm:"not null;default:0"`
}
