package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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
