package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FakeNewrelic is a NewRelic like type.
type FakeNewrelic struct {
	Name string
}

// NewRelicTracer creates a newrelic instance.
func NewRelicTracer(name string) FakeNewrelic {
	return FakeNewrelic{
		Name: fmt.Sprintf("trace %s", name),
	}
}

// Trace logs a tracing message on stdout.
func (n *FakeNewrelic) Trace() {
	log.Printf(n.Name)
}

func userHandler(delayMs time.Duration) http.HandlerFunc {

	nr := NewRelicTracer("users")

	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		nr.Trace()

		userID := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.Println("userID is not a number")
			w.WriteHeader(400)
		}
		user := response{ID: id, Name: "user" + userID}

		time.Sleep(delayMs * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// DB represents a Database interface.
type DB interface {
	SaveUser(u User)
}

// MySQL represents a fake DB implementation.
type MySQL struct{}

// SaveUser is a fake function to create a user.
func (m MySQL) SaveUser(u User) {
	// DB save
}

// User data.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func saveUserHandler(repository DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var msg User
		err = json.Unmarshal(b, &msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		repository.SaveUser(msg)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}
