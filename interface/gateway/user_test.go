// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package gateway

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/interface/iinfra/mock_iinfra"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGatewayFindByEmail(t *testing.T) {
	const query = "SELECT id, name, email FROM users WHERE email = ?"
	const fakeName = "fake name"
	const fakeEmail = "fake@email.com"
	fakeError := errors.New("fake error")

	t.Run("should return an error if the query results in an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer db.Close()

		mock.ExpectQuery(query).WithArgs(fakeEmail).WillReturnError(fakeError)

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any())

		database := mock_iinfra.NewMockDatabase(ctrl)
		database.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, query string, args ...interface{}) (*sql.Rows, error) {
				return db.Query(query, args...)
			})

		g := NewUserGateway(database, logger)
		_, err = g.FindByEmail(context.Background(), fakeEmail)
		assert.EqualError(t, err, fakeError.Error())
	})

	t.Run("should return ErrCreateUserNotFound when query return no results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "email"})
		mock.ExpectQuery(query).WithArgs(fakeEmail).WillReturnRows(rows)

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

		database := mock_iinfra.NewMockDatabase(ctrl)
		database.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, query string, args ...interface{}) (*sql.Rows, error) {
				return db.Query(query, args...)
			})

		g := NewUserGateway(database, logger)
		_, err = g.FindByEmail(context.Background(), fakeEmail)
		assert.EqualError(t, err, businesserr.ErrCreateUserNotFound.Error())
	})

	t.Run("should return an error if occur an error when scanning the result query", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "email"})
		rows.AddRow("invalid id type", fakeName, fakeEmail)
		mock.ExpectQuery(query).WithArgs(fakeEmail).WillReturnRows(rows)

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any())

		database := mock_iinfra.NewMockDatabase(ctrl)
		database.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, query string, args ...interface{}) (*sql.Rows, error) {
				return db.Query(query, args...)
			})

		g := NewUserGateway(database, logger)
		_, err = g.FindByEmail(context.Background(), fakeEmail)
		assert.EqualError(t, err, "sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid id type\") to a int64: invalid syntax")
	})

	t.Run("should return the first user found if the query results any user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "email"})
		rows.AddRow(1, fakeName, fakeEmail)
		mock.ExpectQuery(query).WithArgs(fakeEmail).WillReturnRows(rows)

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

		database := mock_iinfra.NewMockDatabase(ctrl)
		database.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, query string, args ...interface{}) (*sql.Rows, error) {
				return db.Query(query, args...)
			})

		g := NewUserGateway(database, logger)
		user, _ := g.FindByEmail(context.Background(), fakeEmail)
		assert.Equal(t, entity.User{
			ID:    1,
			Name:  fakeName,
			Email: fakeEmail,
		}, user)
	})
}
