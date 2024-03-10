package repositories_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"strconv"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/carlosgonzalez/go-bundled/internal/models"
	"github.com/carlosgonzalez/go-bundled/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository repositories.UserRepository
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, mock, err := sqlmock.New()
	s.mock = mock
	require.NoError(s.T(), err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.DB = DB

	require.NoError(s.T(), err)

	s.repository = repositories.NewUserRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_GetUser() {
	var (
		id   = "1"
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	err, res := s.repository.GetUser("1")

	require.NoError(s.T(), err)

	assert.Equal(s.T(), name, res.Name)

	assert.Equal(s.T(), id, strconv.FormatUint(uint64(res.ID), 10))

}

func (s *Suite) Test_GetAllUsers() {
	var (
		id    = "1"
		name  = "test-name"
		id2   = "2"
		name2 = "test-name2"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(id, name).
				AddRow(id2, name2),
		)

	err, res := s.repository.GetAllUsers()

	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(res))

}

func (s *Suite) Test_CreateUser() {
	var name = "test-name"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(AnyTime{}, AnyTime{}, nil, name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	user := &models.User{
		Name: name,
	}
	err := s.repository.CreateUser(user)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_UpdateUser() {

	oldUser := &models.User{
		Name: "test-name",
	}
	oldUser.ID = 1

	newUser := &models.User{
		Name: "test-name2",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "updated_at"=$1,"name"=$2 WHERE "users"."deleted_at" IS NULL AND "id" = $3`)).
		WithArgs(AnyTime{}, newUser.Name, oldUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err, updatedUser := s.repository.UpdateUser(oldUser, newUser)

	require.NoError(s.T(), err)

	assert.Equal(s.T(), updatedUser.Name, newUser.Name)
}

func (s *Suite) Test_DeleteUser() {

	oldUser := &models.User{
		Name: "test-name",
	}
	oldUser.ID = 1

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."id" = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(AnyTime{}, oldUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.DeleteUser(oldUser)

	require.NoError(s.T(), err)
}
