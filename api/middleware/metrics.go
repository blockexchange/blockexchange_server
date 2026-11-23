package middleware

import (
	"blockexchange/metrics"
	"bufio"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// allows websocket upgrades through the middleware
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw.statusCode = http.StatusSwitchingProtocols
	return http.NewResponseController(rw.ResponseWriter).Hijack()
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// keeps the label values bounded
func metricsMethod(m string) string {
	switch m {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
		http.MethodDelete, http.MethodPatch, http.MethodOptions:
		return m
	default:
		return "other"
	}
}

// path template of the matched route ("/api/schema/{schema_uid}"), never the raw path
func metricsRoute(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route == nil {
		return "unmatched"
	}
	tpl, err := route.GetPathTemplate()
	if err != nil {
		return "unknown"
	}
	return tpl
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		method := metricsMethod(r.Method)
		route := metricsRoute(r)
		metrics.HttpRequests.WithLabelValues(method, route, strconv.Itoa(rw.statusCode)).Inc()

		// websocket handlers return when the connection closes, that's not a response time
		if rw.statusCode != http.StatusSwitchingProtocols {
			metrics.HttpRequestDuration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
		}
	})
}
