package main

import (
	"log"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserStatusSync(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()
	userID := "2"

	start := time.Now()

	result, _ := GetUserStatusSync(srv.URL, userID)

	elapsed := time.Since(start)
	log.Printf("GetUserStatusSync took %s\n", elapsed)
	log.Printf("%+v\n", result)

	if result.ID != userID {
		t.Errorf("unspected result, want %s, got: %s", userID, result.ID)
	}
}

func TestGetUserStatusAsyncWaitgroups(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()
	userID := "2"

	start := time.Now()

	result, _ := GetUserStatusAsyncWaitGroup(srv.URL, userID)

	elapsed := time.Since(start)
	log.Printf("GetUserStatusAsyncWaitGroup took %s\n", elapsed)
	log.Printf("%+v\n", result)

	if result.ID != userID {
		t.Errorf("unspected result, want %s, got: %s", userID, result.ID)
	}
}

func TestGetUserStatusAsyncChannels(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()
	userID := "2"

	start := time.Now()

	result, _ := GetUserStatusAsyncChannels(srv.URL, userID)

	elapsed := time.Since(start)
	log.Printf("GetUserStatusAsyncChannels took %s\n", elapsed)
	log.Printf("%+v\n", result)

	if result.ID != userID {
		t.Errorf("unspected result, want %s, got: %s", userID, result.ID)
	}
}
