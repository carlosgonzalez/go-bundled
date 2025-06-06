package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carlosgonzalez/go-bundled/internal/handlers"
	"github.com/carlosgonzalez/go-bundled/internal/models"
	"github.com/carlosgonzalez/go-bundled/mocks"
	"github.com/carlosgonzalez/go-bundled/pkg/validators"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	userJSON = `{"name": "FooBar"}`
)

func TestUserHandler_CreateUserSucceedsForValidPayload(t *testing.T) {
	c, rec := getContextForRoute("/users", http.MethodPost)

	repo := &mocks.MockUserRepositoryInterface{}
	repo.On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(nil).
		Once()

	handler := handlers.NewUserHandler(repo)

	// Assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		expectedUser := &models.User{}
		err := json.Unmarshal([]byte(userJSON), expectedUser)
		assert.NoError(t, err)

		actualUser := &models.User{}
		err = json.Unmarshal(rec.Body.Bytes(), actualUser)
		assert.NoError(t, err)

		assert.Equal(t, expectedUser, actualUser)
	}

}

func TestUserHandler_CreateUserFailsForInvalidPayload(t *testing.T) {
	c, rec := getContextForRoute("/users", http.MethodPost)

	repo := &mocks.MockUserRepositoryInterface{}
	repo.On("CreateUser", mock.AnythingOfType("*models.User")).
		Return(errors.New("something went wrong")).
		Once()

	handler := handlers.NewUserHandler(repo)

	// Assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

}

func getContextForRoute(endpoint string, method string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = validators.NewCustomValidator()
	req := httptest.NewRequest(method, endpoint, strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
