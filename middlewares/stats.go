package middlewares

import (
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}

	LogInfo struct {
		Method              string `json:"method"`
		Path                string `json:"path"`
		LastRequestDuration int64  `json:"last_request_duration"`
	}
)

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: map[string]int{},
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		before := time.Now()

		logInfo := LogInfo{
			Path: c.Path(),
		}

		c.Response().Before(func() {
			logInfo.LastRequestDuration = time.Since(before).Milliseconds()
			logInfo.Method = c.Request().Method
		})

		if err := next(c); err != nil {
			c.Error(err)
		}

		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++

		c.Echo().Logger.Info(logInfo)

		return nil
	}
}
