package server

import (
	"fmt"
	"net/http"
)

func GetOne(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getone" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("search")

	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u, err := base.findUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := fmt.Fprintf(w, "search result:%v", u); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func SaveOne(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/saveone" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("user")

	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u := user{
		Username: username,
		Email:    "test@mail.ru",
	}

	err := base.saveUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := fmt.Fprintf(w, "saved:%v", u.Username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getall" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	res := base.showAll()

	if _, err := fmt.Fprintf(w, "result:%v", res); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
