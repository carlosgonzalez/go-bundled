package repositories_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/carlosgonzalez/learning-go/internal/models"
	"github.com/carlosgonzalez/learning-go/internal/repositories"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository repositories.UserRepository
	user       *models.User
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, mock, err := sqlmock.New()
	s.mock = mock
	require.NoError(s.T(), err)

	// DB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	DB, err := gorm.Open(sqlite.Dialector{
		Conn:       db,
		DriverName: "sqlite",
	}, &gorm.Config{})
	s.DB = DB

	require.NoError(s.T(), err)

	s.repository = *repositories.NewUserRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_repository_Get() {
	var (
		id   = "1"
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "person" WHERE (id = $1)`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name))

	err, res := s.repository.GetUser("1")

	require.NoError(s.T(), err)
	fmt.Println(res)
}
