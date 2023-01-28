package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (m *AuthMiddleware) LoginHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		username, password := extractCredentials(r)

		u, err := m.Login(username, password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid credentials"))
			return
		}

		accessToken, refreshToken, err := m.GenerateAndSendTokens(u)
		if err != nil {
			log.Printf("GenerateAndSendTokens: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error proceeding tokens"))
			return
		}

		t := &token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		reply, err := json.Marshal(t)
		if err != nil {
			log.Printf("Login: Marshal token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error proceeding tokens"))
			return
		}

		w.Write(reply)
	}
}

func (m *AuthMiddleware) UpdateHanlder() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		refreshToken := string(r.FormValue("refresh"))

		accessToken, refreshToken, err := m.UpdateToken(refreshToken)

		if err != nil {
			log.Printf("GenerateAndSendTokens: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error proceeding tokens"))
			return
		}

		t := &token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		reply, err := json.Marshal(t)
		if err != nil {
			log.Printf("Update: Marshal token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error proceeding tokens"))
			return
		}

		w.Write(reply)
	}
}

func extractCredentials(r *http.Request) (login, pass string) {
	return string(r.FormValue("username")), string(r.FormValue("password"))
}
