package database

import (
	"github.com/dougefr/go-clean-code/infrastructure/database/entity"
	"github.com/jinzhu/gorm"
)

type database struct {
	db *gorm.DB
}

// Database ...
type Database interface {
	Create(value interface{}) error
	First(out interface{}, where ...interface{}) error
}

// NewDatabase ...
func NewDatabase() (Database, error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.User{})

	return database{
		db: db,
	}, nil
}

// Create ...
func (d database) Create(value interface{}) error {
	return d.db.Create(value).Error
}

// First ...
func (d database) First(out interface{}, where ...interface{}) error {
	return d.db.First(out, where).Error
}
