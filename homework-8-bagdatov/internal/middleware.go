package internal

import (
	"context"
	"net/http"
	"time"
)

type ctxKey int8

const ctxKeyUser ctxKey = iota

type handler func(w http.ResponseWriter, r *http.Request)

func (ub *UserBase) Middleware(next handler) http.HandlerFunc {
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

		user, err := ub.FindUser(username, password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
		next(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}
