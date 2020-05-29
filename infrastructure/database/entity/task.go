package entity

import "github.com/jinzhu/gorm"

// Task ...
type Task struct {
	gorm.Model
	Text     string
	Done     bool
	AssignTo User
}
