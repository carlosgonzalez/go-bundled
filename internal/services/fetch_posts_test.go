package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchUrl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/posts/1", r.URL.String())
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"userId": 1,
				"id":     1,
				"title":  "Test Title",
				"body":   "Test Body",
			}
			err := json.NewEncoder(w).Encode(response)
			assert.NoError(t, err)
		}))
		defer server.Close()

		post, err := FetchUrl(server.Client(), server.URL, "posts", 1)

		assert.NoError(t, err)
		assert.Equal(t, "Test Title", post.Title)
		assert.Equal(t, "Test Body", post.Body)
		assert.Equal(t, uint(0), post.ID, "ID should be reset to 0")
	})

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close() // Close the server immediately

		_, err := FetchUrl(server.Client(), server.URL, "posts", 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error making http request")
	})

	t.Run("bad status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, err := FetchUrl(server.Client(), server.URL, "posts", 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("bad json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id":1, "title":"test",`)) // malformed json
			assert.NoError(t, err)
		}))
		defer server.Close()

		_, err := FetchUrl(server.Client(), server.URL, "posts", 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling json")
	})
}

func TestFetcher(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"userId": 1,
				"id":     1,
				"title":  "Test Title",
				"body":   "Test Body",
			}
			err := json.NewEncoder(w).Encode(response)
			assert.NoError(t, err)
		}))
		defer server.Close()

		posts, err := Fetcher(server.URL, "posts", 5)

		assert.NoError(t, err)
		assert.Len(t, posts, 5)
		for _, post := range posts {
			assert.Equal(t, "Test Title", post.Title)
			assert.Equal(t, "Test Body", post.Body)
		}
	})

	t.Run("error from FetchUrl", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		posts, err := Fetcher(server.URL, "posts", 5)

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})
}
