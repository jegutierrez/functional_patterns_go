package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	err := http.ListenAndServe(":8080", handler())
	if err != nil {
		log.Fatal(err)
	}
}

func handler() http.Handler {
	srv := http.NewServeMux()
	srv.HandleFunc("/users/", userHandler)
	srv.HandleFunc("/balance/", balanceHandler)
	srv.HandleFunc("/user-debts/", debtsHandler)
	log.Println("server listening connections")
	return srv
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	user := map[string]string{"id": userID, "name": "user" + userID}

	// delay
	time.Sleep(150 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/balance/")
	amount := fmt.Sprint(rand.Intn(100), ",", rand.Intn(100))
	balance := map[string]string{"user_id": userID, "amount": amount}

	// delay
	time.Sleep(350 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

func debtsHandler(w http.ResponseWriter, r *http.Request) {
	debts := []struct {
		ID     string `json:"id"`
		Reason string `json:"reason"`
		Amount string `json:"amount"`
	}{
		{ID: "14", Reason: "chargeback", Amount: "71.0"},
		{ID: "37", Reason: "chargeback", Amount: "15.5"},
		{ID: "51", Reason: "chargeback", Amount: "43.0"},
	}
	// delay
	time.Sleep(350 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debts)
}
