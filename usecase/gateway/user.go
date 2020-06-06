package gateway

import (
	"context"
	"github.com/dougefr/go-clean-arch/entity"
)

// User ...
type User interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}
