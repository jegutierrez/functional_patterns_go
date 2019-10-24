package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockDB struct {
	MockSaveUserFn func(User)
}

func (m MockDB) SaveUser(u User) {
	m.MockSaveUserFn(u)
}

func helperMockDB(t *testing.T) func(User) {
	t.Helper()

	return func(u User) {
		if u.ID != 0 {
			t.Errorf("user ID must not be preset")
		}
	}
}

func TestHttpHandler(t *testing.T) {
	body := strings.NewReader(`{"name": "john"}`)
	req, err := http.NewRequest("POST", "/users", body)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	saveFn := helperMockDB(t)
	mockDB := MockDB{
		MockSaveUserFn: saveFn,
	}

	saveUser := saveUserHandler(mockDB)

	handler := http.HandlerFunc(saveUser)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
