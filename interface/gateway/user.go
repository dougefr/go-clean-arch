// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dougefr/go-clean-arch/entity"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	"github.com/dougefr/go-clean-arch/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/usecase/igateway"

	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// default error when query execution fails
const errorExecutingQuery = "error when executing query: %v"

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
	startTime := time.Now()
	u.logger.Debug(ctx, "starting find by email method")

	var rows *sql.Rows
	rows, err = u.db.Query(ctx, "SELECT id, name, email FROM users WHERE email = ?", email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf(errorExecutingQuery, err), iinfra.LogAttrs{"email": email})
		return
	}

	if rows.Next() {
		// get just the first line
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when scanning query result: %v", err),
				iinfra.LogAttrs{"email": email})
			return
		}
	} else {
		// will return an error if the user does not exists
		err = businesserr.ErrCreateUserNotFound
	}

	u.logger.Debug(ctx, "ending find by email method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}

// Create ...
func (u userGateway) Create(ctx context.Context, user entity.User) (userCreated entity.User, err error) {
	startTime := time.Now()
	u.logger.Debug(ctx, "starting create user method")

	result, err := u.db.Exec(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf(errorExecutingQuery, err), iinfra.LogAttrs{"user": user})
		return
	}

	// get the ID of the user that was created
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
	startTime := time.Now()
	u.logger.Debug(ctx, "starting find all users method")

	var rows *sql.Rows
	rows, err = u.db.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf(errorExecutingQuery, err))
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

	u.logger.Debug(ctx, "ending find all users method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}
