package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/carlosgonzalez/go-bundled/internal/models"
)

func Fetcher(resource string, totalRecords int) ([]models.Post, error) {

	c := make(chan models.Post, totalRecords)
	wg := sync.WaitGroup{}

	wg.Add(totalRecords)
	for i := 1; i <= totalRecords; i++ {
		go fetchUrl(resource, c, &wg)
	}

	wg.Wait()
	close(c)

	var posts []models.Post
	for post := range c {
		posts = append(posts, post)
	}

	return posts, nil
}

func fetchUrl(resource string, c chan models.Post, wg *sync.WaitGroup) {
	postNumber := rand.Intn(100)
	requestURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/%s/%d", resource, postNumber)

	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("error making http request: %s\n", err)
		os.Exit(1)
	}

	var post models.Post

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &post); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	//to not mess with the autoincrementing ID from the DB
	post.ID = 0

	c <- post
	wg.Done()
}
