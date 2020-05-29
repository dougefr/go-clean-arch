package repository

import (
	"github.com/dougefr/go-clean-code/entity"
	"github.com/dougefr/go-clean-code/infrastructure/database"
	. "github.com/dougefr/go-clean-code/infrastructure/database/entity"
	"github.com/dougefr/go-clean-code/usecase"
	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type userRepo struct {
	db database.Database
}

// NewUserRepo ...
func NewUserRepo(db database.Database) usecase.UserRepo {
	return userRepo{
		db: db,
	}
}

// FindByEmail ...
func (u userRepo) FindByEmail(email string) (entity.User, error) {
	var user User
	err := u.db.First(&user, "email = ?", email)

	return entity.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, err
}

// CreateUser ...
func (u userRepo) CreateUser(user entity.User) (entity.User, error) {
	err := u.db.Create(&User{
		Name:  user.Name,
		Email: user.Email,
	})

	return entity.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, err
}
