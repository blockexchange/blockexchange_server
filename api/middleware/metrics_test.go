package middleware_test

import (
	"blockexchange/api/middleware"
	"blockexchange/metrics"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware(t *testing.T) {
	r := mux.NewRouter()
	r.Use(middleware.PrometheusMiddleware)
	r.HandleFunc("/api/schema/{schema_uid}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	})

	requests := metrics.HttpRequests.WithLabelValues("GET", "/api/schema/{schema_uid}", "404")
	other := metrics.HttpRequests.WithLabelValues("other", "/api/schema/{schema_uid}", "404")
	before, before_other := testutil.ToFloat64(requests), testutil.ToFloat64(other)
	histograms := testutil.CollectAndCount(metrics.HttpRequestDuration)

	// the path template is the label, not the raw path
	for _, uid := range []string{"a", "b"} {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/schema/"+uid, nil))
	}
	assert.Equal(t, before+2, testutil.ToFloat64(requests))

	// unknown methods are folded together
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("FOO", "/api/schema/x", nil))
	assert.Equal(t, before_other+1, testutil.ToFloat64(other))

	// websocket connections are counted but not timed
	ws := metrics.HttpRequests.WithLabelValues("GET", "/api/ws", "101")
	before_ws := testutil.ToFloat64(ws)
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/ws", nil))
	assert.Equal(t, before_ws+1, testutil.ToFloat64(ws))
	// series: GET and "other" on the schema route, none for the websocket route
	assert.Equal(t, histograms+2, testutil.CollectAndCount(metrics.HttpRequestDuration))
}
