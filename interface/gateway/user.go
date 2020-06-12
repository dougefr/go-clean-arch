package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/igateway"
	usecase2 "github.com/dougefr/go-clean-arch/core/usecase"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	"time"

	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type userGateway struct {
	db     iinfra.Database
	logger iinfra.LogProvider
}

// NewUserGateway ...
func NewUserGateway(db iinfra.Database, logger iinfra.LogProvider) igateway.User {
	return userGateway{
		db:     db,
		logger: logger,
	}
}

// FindByEmail ...
func (u userGateway) FindByEmail(ctx context.Context, email string) (user entity.User, err error) {
	var rows *sql.Rows
	rows, err = u.db.Query(ctx, "SELECT id, name, email FROM users WHERE email = ?", email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing query: %v", err), iinfra.LogAttrs{"email": email})
		return
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when scanning query result: %v", err), iinfra.LogAttrs{"email": email})
			return
		}
	} else {
		err = usecase2.ErrCreateUserNotFound
	}

	return
}

// CreateUser ...
func (u userGateway) Create(ctx context.Context, user entity.User) (userCreated entity.User, err error) {
	startTime := time.Now()
	u.logger.Debug(ctx, "starting create user method")

	result, err := u.db.Exec(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing query: %v", err), iinfra.LogAttrs{"user": user})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when getting last insert ID: %v", err), iinfra.LogAttrs{"user": user})
		return
	}

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return entity.User{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, err
}

// FindAll ...
func (u userGateway) FindAll(ctx context.Context) (users []entity.User, err error) {
	var rows *sql.Rows
	rows, err = u.db.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing query: %v", err))
		return
	}

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when scanning query result: %v", err))
			return
		}

		users = append(users, user)
	}

	return
}
