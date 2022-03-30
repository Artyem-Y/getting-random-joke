package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	fmt.Fprintf(w, "handler\n")
}

func getRandomResult() {
	respPerson, err := http.Get("https://names.mcquay.me/api/v0/")
	if err != nil {
		log.Fatal("err: ", err)
	}

	log.Println("Fetching a joke")

	bodyPerson, err := ioutil.ReadAll(respPerson.Body)
	if err != nil {
		log.Fatal("err: ", err)
	}

	respJoke, err := http.Get("http://api.icndb.com/jokes/random?firstName=John&lastName=Doe&limitTo=nerdy")
	if err != nil {
		log.Fatal("err: ", err)
	}

	bodyJoke, err := ioutil.ReadAll(respJoke.Body)
	if err != nil {
		log.Fatal("err: ", err)
	}

	var person Person
	var jokeResp JokeResp

	err = json.Unmarshal(bodyPerson, &person)
	if err != nil {
		log.Fatal("err: ", err)
	}

	err = json.Unmarshal(bodyJoke, &jokeResp)
	if err != nil {
		log.Fatal("err: ", err)
	}

	log.Println(person.FirstName + " " + person.LastName + "'s " + jokeResp.Value.Joke)

	os.Exit(0)
}

func main() {
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

	getRandomResult()

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start with error %v", err.Error())
	}
}
