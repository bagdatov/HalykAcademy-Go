package middleware

import (
	"context"
	"net/http"
	"time"
	"web/domain"

	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	domain.AuthUseCase
}

type ctxKey int8

type handler func(w http.ResponseWriter, r *http.Request)

const CtxKeyUser ctxKey = iota

func InitMiddleware(uc domain.AuthUseCase) AuthMiddleware {
	return AuthMiddleware{uc}
}

func (m *AuthMiddleware) CheckAuthMiddleware(next handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		timer := time.Now()
		w.Header().Set("Trailer", "Execution")

		defer func() {
			// if you don't cover this defer with func() it will evaluate time incorrectly
			w.Header().Set("Execution", time.Since(timer).String())
		}()

		header := r.Header.Get("Authorization")

		token, err := m.ExtractToken(header)
		if err != nil {
			log.Printf("Extract token error: %v", err)
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		user, err := m.ParseToken(token, true)
		if err != nil {
			if err == domain.ErrExpiredToken {
				http.Redirect(w, r, "/update", http.StatusMovedPermanently)
				return
			}
			log.Printf("Parse access token error: %v", err)
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, user)))
	}
}
