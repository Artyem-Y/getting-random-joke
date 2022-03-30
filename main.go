package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpHandler struct {
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Joke struct {
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

type JokeResp struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	log.Println("handler")
}

func getRandomJoke() (string, error) {
	respPerson, err := http.Get("https://names.mcquay.me/api/v0/")
	if err != nil {
		return "", err
	}

	log.Println("Fetching a joke")

	respJoke, err := http.Get("http://api.icndb.com/jokes/random?firstName=John&lastName=Doe&limitTo=nerdy")
	if err != nil {
		return "", err
	}

	var person Person
	var jokeResp JokeResp

	err = json.NewDecoder(respPerson.Body).Decode(&person)
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(respJoke.Body).Decode(&jokeResp)
	if err != nil {
		return "", err
	}

	return person.FirstName + " " + person.LastName + "'s " + jokeResp.Value.Joke, nil
}

func main() {
	start := time.Now()
	defer func() {
		log.Println("Execution Time: ", time.Since(start))
	}()

	// Define a serveMux to handle routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Custom http server.
	s := &http.Server{
		Addr: ":5000",
		// Wrap the servemux with the limit middleware.
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var result string
	var err error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		result, err = getRandomJoke()
		if err != nil {
			log.Fatal("err: ", err)
			return
		}

		log.Println(result)
		wg.Done()
	}()

	wg.Wait()

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start with error %v", err.Error())
	}
}
