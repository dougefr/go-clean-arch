/*
 * Go Clean Architecture
 * Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the Software
 * without restriction, including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
 * to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 * PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
 * FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 *
 * Source available at: https://github.com/dougefr/go-clean-arch
 */

package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/core/usecase/igateway"
	"github.com/dougefr/go-clean-arch/interface/iinfra"

	// sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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
	var rows *sql.Rows
	rows, err = u.db.Query(ctx, "SELECT id, name, email FROM users WHERE email = ?", email)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf(errorExecutingQuery, err), iinfra.LogAttrs{"email": email})
		return
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when scanning query result: %v", err), iinfra.LogAttrs{"email": email})
			return
		}
	} else {
		err = businesserr.ErrCreateUserNotFound
	}

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

	return
}
