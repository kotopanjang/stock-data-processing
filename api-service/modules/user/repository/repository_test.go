package repository_test

import (
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"stock-data-processing/api-service/modules/user/repository"
)

func Test_repository_GetOne(t *testing.T) {

	mockDb, mockSql, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDb.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDb,
		PreferSimpleProtocol: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)
	if conn == nil {
		log.Fatal("Failed to open connection to DB: conn is nil")
	}

	rows := sqlmock.NewRows([]string{"id", "name", "age", "created_at", "deleted_at", "updated_at"}).
		AddRow(1, "test-1", 1, time.Now(), time.Now(), time.Now())
	mockSql.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM users ORDER BY id LIMIT 1`)).
		WillReturnRows(rows)

	repo := repository.NewRepository(conn)
	_, err = repo.GetOne(1)

	assert.NoError(t, err)

}

func Test_repository_GetAll(t *testing.T) {

	mockDb, mockSql, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDb.Close()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDb,
		PreferSimpleProtocol: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)
	if conn == nil {
		log.Fatal("Failed to open connection to DB: conn is nil")
	}

	rows := sqlmock.NewRows([]string{"id", "name", "age", "created_at", "deleted_at", "updated_at"}).
		AddRow(1, "test-1", 12, time.Now(), time.Now(), time.Now()).
		AddRow(2, "test-2", 13, time.Now(), time.Now(), time.Now())
	mockSql.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM users`)).
		WillReturnRows(rows)

	repo := repository.NewRepository(conn)
	_, _, err = repo.GetList(2, 0)

	assert.NoError(t, err)

}

// func Test_repository_Create(t *testing.T) {

// 	mockDb, mockSql, err := sqlmock.New()
// 	assert.Nil(t, err)
// 	defer mockDb.Close()

// 	dialector := postgres.New(postgres.Config{
// 		DSN:                  "sqlmock_db_0",
// 		DriverName:           "postgres",
// 		Conn:                 mockDb,
// 		PreferSimpleProtocol: true,
// 	})

// 	conn, err := gorm.Open(dialector, &gorm.Config{})
// 	require.NoError(t, err)
// 	if conn == nil {
// 		log.Fatal("Failed to open connection to DB: conn is nil")
// 	}

// 	var (
// 		// id   = 123
// 		name = "test-name"
// 		age  = 13
// 	)

// 	// mockSql.MatchExpectationsInOrder(false)
// 	mockSql.ExpectBegin()
// 	mockSql.ExpectExec(regexp.QuoteMeta(
// 		`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","age") VALUES ($1,$2,$3,$4,$5) RETURNING "users"."id"`)).
// 		WithArgs(time.Now(), time.Now(), time.Now(), name, age).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mockSql.ExpectCommit()

// 	// repoMock := new(repoMock.Repository)
// 	// repoMock.On("InsertOne", mock.Anything).Return(nil)
// 	// repoMock.AssertExpectations(t)

// 	repo := repository.NewRepository(conn)
// 	err = repo.InsertOne(entities.User{
// 		Name: "test",
// 		Age:  13,
// 	})

// 	assert.NoError(t, err)

// }
