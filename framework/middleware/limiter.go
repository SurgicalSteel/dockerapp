package middleware

import "net/http"

const Limit = 100

func LimitRate(h http.Handler) http.Handler {
	limit := make(chan bool, Limit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Limit == len(limit) {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			limit <- true
			defer func() { <-limit }()
			h.ServeHTTP(w, r)
		}
	})
}
