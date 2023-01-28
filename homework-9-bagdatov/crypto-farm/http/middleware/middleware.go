package middleware

import (
	"context"
	"net/http"
	"time"
	"web/domain"
)

type MyMiddleware struct {
	domain.UseCase
}

type ctxKey int8

type handler func(w http.ResponseWriter, r *http.Request)

const CtxKeyUser ctxKey = iota

func InitMiddleware(uc domain.UseCase) *MyMiddleware {
	return &MyMiddleware{uc}
}

func (m *MyMiddleware) Middleware(next handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		timer := time.Now()
		w.Header().Set("Trailer", "Execution")

		defer func() {
			// if you don't cover this defer with func() it will evaluate time incorrectly
			w.Header().Set("Execution", time.Since(timer).String())
		}()

		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := m.FindUser(r.Context(), username, password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, user)))
	}
}
