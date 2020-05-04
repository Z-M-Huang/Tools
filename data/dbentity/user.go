package dbentity

import "github.com/jinzhu/gorm"

//User database entity
type User struct {
	gorm.Model
	Username string
	Password string
	Email    string `gorm:"primary_key"`
	GoogleID string
}
