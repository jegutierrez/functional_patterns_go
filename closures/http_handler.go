package main

import (
	"encoding/json"
	"fmt"
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
