package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"

	"github.com/carlosgonzalez/go-bundled/internal/models"
)

var Fetcher = func(baseURL string, resource string, totalRecords int) ([]models.Post, error) {

	postsChan := make(chan models.Post, totalRecords)
	errChan := make(chan error, totalRecords)
	wg := sync.WaitGroup{}
	client := &http.Client{}

	wg.Add(totalRecords)
	for i := 1; i <= totalRecords; i++ {
		go func() {
			defer wg.Done()
			postNumber := rand.Intn(100)
			if postNumber == 0 {
				postNumber = 1
			}
			post, err := FetchUrl(client, baseURL, resource, postNumber)
			if err != nil {
				errChan <- err
				return
			}
			postsChan <- post
		}()
	}

	wg.Wait()
	close(postsChan)
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	var posts []models.Post
	for post := range postsChan {
		posts = append(posts, post)
	}

	return posts, nil
}

// FetchUrl fetches a single post from the given resource. It is exported for testing purposes.
var FetchUrl = func(client *http.Client, baseURL string, resource string, postNumber int) (models.Post, error) {
	var post models.Post
	requestURL := fmt.Sprintf("%s/%s/%d", baseURL, resource, postNumber)

	resp, err := client.Get(requestURL)
	if err != nil {
		return post, fmt.Errorf("error making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return post, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, &post); err != nil {
		return post, fmt.Errorf("error unmarshaling json: %w", err)
	}

	//to not mess with the autoincrementing ID from the DB
	post.ID = 0

	return post, nil
}
