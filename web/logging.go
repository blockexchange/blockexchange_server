package web

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		logrus.WithFields(logrus.Fields{
			"origin": r.Header.Get("X-Forwarded-For"),
			"method": r.Method,
			"path":   r.URL.Path,
			"status": rw.statusCode,
		}).Info("HTTP Request")
	})
}
