package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dougefr/go-clean-arch/entity"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	"github.com/dougefr/go-clean-arch/usecase"
	"github.com/dougefr/go-clean-arch/usecase/gateway"
	"time"

	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type userRepo struct {
	db     iinfra.Database
	logger iinfra.LogProvider
}

// NewUserRepo ...
func NewUserRepo(db iinfra.Database, logger iinfra.LogProvider) gateway.User {
	return userRepo{
		db:     db,
		logger: logger,
	}
}

// FindByEmail ...
func (u userRepo) FindByEmail(ctx context.Context, email string) (user entity.User, err error) {
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
		err = usecase.ErrCreateUserNotFound
	}

	return
}

// CreateUser ...
func (u userRepo) CreateUser(ctx context.Context, user entity.User) (userCreated entity.User, err error) {
	startTime := time.Now()
	u.logger.Debug(ctx, "starting create user method")

	_, err = u.db.Exec(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing query: %v", err), iinfra.LogAttrs{"user": user})
		return
	}

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return entity.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, err
}

// FindAll ...
func (u userRepo) FindAll(ctx context.Context) (users []entity.User, err error) {
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
