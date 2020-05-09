package dbentity

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//User database entity
type User struct {
	gorm.Model
	Username  string `gorm:"unique_index"`
	Password  string
	Email     string `gorm:"primary_key"`
	GoogleID  string
	LikedApps pq.StringArray `gorm:"type:text[]"`
}
