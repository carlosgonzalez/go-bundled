package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carlosgonzalez/learning-go/internal/handlers"
	"github.com/carlosgonzalez/learning-go/internal/models"
	"github.com/carlosgonzalez/learning-go/mocks"
	"github.com/carlosgonzalez/learning-go/pkg/validators"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var userJSON = `{"name": "FooBar"}`

func TestUserHandler_CreateUser(t *testing.T) {

	//Setup
	e := echo.New()
	e.Validator = validators.NewCustomValidator()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := &mocks.UserRepositoryInterface{}
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
