package middleware

import (
	"github.com/febytanzil/dockerapp/framework/config"
	"net/http"
)

func LimitRate(h http.HandlerFunc) http.HandlerFunc {
	limit := config.Get().App.Limit
	limitChan := make(chan bool, limit)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limit == len(limitChan) {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			limitChan <- true
			defer func() { <-limitChan }()
			h.ServeHTTP(w, r)
		}
	})
}
