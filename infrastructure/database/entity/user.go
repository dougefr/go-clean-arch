package entity

import "github.com/jinzhu/gorm"

// User ...
type User struct {
	gorm.Model
	Name  string
	Email string
}
