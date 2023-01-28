package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetOne(t *testing.T) {
	// check existing user
	r := httptest.NewRequest(http.MethodGet, "/getone?search=alibi", nil)
	w := httptest.NewRecorder()
	GetOne(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Response code failed")
	}

	// check non-existing user
	r = httptest.NewRequest(http.MethodGet, "/getone?search=nonExistingUser", nil)
	w = httptest.NewRecorder()
	GetOne(w, r)

	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Response code failed")
	}

	// check incorrect method
	r = httptest.NewRequest(http.MethodPost, "/getone", nil)
	w = httptest.NewRecorder()
	GetOne(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Response code failed")
	}

	// check empty request
	r = httptest.NewRequest(http.MethodGet, "/getone", nil)
	w = httptest.NewRecorder()
	GetOne(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Response code failed")
	}

}

func TestSaveOne(t *testing.T) {
	// create new user
	r := httptest.NewRequest(http.MethodPost, "/saveone", strings.NewReader("user=testUser"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	SaveOne(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Response code failed")
	}

	// create existing user
	r = httptest.NewRequest(http.MethodPost, "/saveone", strings.NewReader("user=alibi"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	SaveOne(w, r)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Response code failed")
	}

	// create invalid user
	r = httptest.NewRequest(http.MethodPost, "/saveone", strings.NewReader("user="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	SaveOne(w, r)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Response code failed")
	}

	// incorrect method
	r = httptest.NewRequest(http.MethodGet, "/saveone", strings.NewReader("user="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	SaveOne(w, r)

	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Response code failed")
	}
}

func TestGetAll(t *testing.T) {
	// show all
	r := httptest.NewRequest(http.MethodGet, "/getall", nil)

	w := httptest.NewRecorder()
	GetAll(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		t.Fatalf("Response code failed")
	}

	// incorrect method
	r = httptest.NewRequest(http.MethodPost, "/getall", nil)
	w = httptest.NewRecorder()
	GetAll(w, r)

	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		fmt.Println(resp.StatusCode)
		t.Fatalf("Response code failed")
	}
}
