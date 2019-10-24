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

func balanceHandler(delayMs time.Duration) http.HandlerFunc {

	nr := NewRelicTracer("balances")

	type response struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		nr.Trace()

		balanceUserID := strings.TrimPrefix(r.URL.Path, "/balance/")
		userID, err := strconv.Atoi(balanceUserID)
		if err != nil {
			log.Println("balanceID is not a number")
			w.WriteHeader(400)
		}
		balance := response{UserID: userID, Amount: 100}

		time.Sleep(delayMs * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(balance)
	}
}
