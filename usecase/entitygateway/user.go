package entitygateway

import "github.com/dougefr/go-clean-code/entity"

// UserRepo ...
type User interface {
	FindByEmail(email string) (entity.User, error)
	CreateUser(user entity.User) (entity.User, error)
}