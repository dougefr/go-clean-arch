package igateway

import (
	"context"

	"github.com/dougefr/go-clean-arch/core/entity"
)

// User ...
type User interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
}
