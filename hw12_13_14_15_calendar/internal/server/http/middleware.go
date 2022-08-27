package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(logg Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		dur := time.Since(start)
		logg.Info(fmt.Sprintf("%v, %v, %v, %v, %v - processed in %v",
			r.RemoteAddr, start, r.Method, r.RequestURI, r.UserAgent(), dur))
	})
}
