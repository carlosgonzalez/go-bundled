package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/carlosgonzalez/go-bundled/internal/models"
	"github.com/carlosgonzalez/go-bundled/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		handler := NewPostHandler(gormDB)

		originalFetcher := services.Fetcher
		defer func() { services.Fetcher = originalFetcher }()
		services.Fetcher = func(baseURL string, resource string, totalRecords int) ([]models.Post, error) {
			return []models.Post{
				{Title: "Test Post 1", Body: "Test Body 1"},
				{Title: "Test Post 2", Body: "Test Body 2"},
			}, nil
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "posts" (.+) VALUES (.+),(.+) RETURNING "id"`).
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Test Post 1", "Test Body 1",
				sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Test Post 2", "Test Body 2",
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
		mock.ExpectCommit()

		err = handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, mock.ExpectationsWereMet())

		var posts []models.Post
		err = json.Unmarshal(rec.Body.Bytes(), &posts)
		assert.NoError(t, err)
		assert.Len(t, posts, 2)
		assert.Equal(t, "Test Post 1", posts[0].Title)
	})

	t.Run("fetcher error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		handler := NewPostHandler(gormDB)

		originalFetcher := services.Fetcher
		defer func() { services.Fetcher = originalFetcher }()
		services.Fetcher = func(baseURL string, resource string, totalRecords int) ([]models.Post, error) {
			return nil, errors.New("fetcher error")
		}

		err = handler.CreatePost(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"error":"unable to fetch posts"}`, rec.Body.String())
	})
}
