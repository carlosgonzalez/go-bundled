package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlosgonzalez/learning-go/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStats_Process(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	stats := middlewares.NewStats()

	// Assert initial state
	assert.Equal(t, uint64(0), stats.RequestCount)
	assert.Equal(t, 0, len(stats.Statuses))

	// Call
	err := stats.Process(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	// Assert after call
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), stats.RequestCount)
	assert.Equal(t, 1, len(stats.Statuses))
	assert.Equal(t, 1, stats.Statuses["200"])
}

func TestStats_ProcessError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	stats := middlewares.NewStats()

	// Assert initial state
	assert.Equal(t, uint64(0), stats.RequestCount)
	assert.Equal(t, 0, len(stats.Statuses))

	// Call
	err := stats.Process(func(c echo.Context) error {
		return echo.ErrNotFound
	})(c)

	// Assert after call
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), stats.RequestCount)
	assert.Equal(t, 1, len(stats.Statuses))
	assert.Equal(t, 1, stats.Statuses["404"])
}
