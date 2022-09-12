package web

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		host := r.Host
		forwHost := r.Header.Get("X-Forwarded-Host")
		if forwHost != "" {
			host = forwHost
		}

		logrus.WithFields(logrus.Fields{
			"host":   host,
			"method": r.Method,
			"path":   r.URL.Path,
			"status": rw.statusCode,
		}).Info("HTTP Request")
	})
}
