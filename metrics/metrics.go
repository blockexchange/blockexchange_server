// prometheus metrics of the blockexchange server, served on /metrics
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// http

var HttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_http_requests_total",
	Help: "Number of handled http requests",
}, []string{"method", "route", "status"})

var HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "bx_http_request_duration_seconds",
	Help:    "Duration of the http requests (without websocket connections)",
	Buckets: []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30},
}, []string{"method", "route"})

// schemaparts

var SchemaPartsUploaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_uploaded",
	Help: "The total number of uploaded schemaparts",
})

var SchemaPartsDownloaded = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_schemaparts_downloaded",
	Help: "The total number of downloaded schemaparts",
})

// background jobs

var JobRuns = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_job_runs_total",
	Help: "Number of background job runs",
}, []string{"name", "result"})

var JobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "bx_job_duration_seconds",
	Help:    "Duration of the background job runs",
	Buckets: []float64{0.1, 0.5, 1, 5, 10, 30, 60, 300, 900},
}, []string{"name"})

var JobLockErrors = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_job_lock_errors_total",
	Help: "Number of failed attempts to acquire the cluster-wide job lock",
}, []string{"name"})

var JobLastSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bx_job_last_success_timestamp_seconds",
	Help: "Unix time of the last successful job run on this node",
}, []string{"name"})

// content stats (updated by the stats job)

var Schemas = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_schemas",
	Help: "Number of schematics",
})

var SchemaParts = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_schemaparts",
	Help: "Number of schematic parts",
})

var Users = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_users",
	Help: "Number of users",
})

var StorageBytes = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_storage_bytes",
	Help: "Total size of all schematics",
})

// wasm packs

var WasmPackRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_pack_requests_total",
	Help: "Wasm pack requests by result (hit: cached, miss: built for this request alone, shared: build served several concurrent requests, error, canceled)",
}, []string{"kind", "result"})

var WasmPackBuildDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "bx_wasm_pack_build_duration_seconds",
	Help:    "Duration of the wasm pack builds",
	Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
}, []string{"kind", "result"})

var WasmPackResolves = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_pack_resolve_total",
	Help: "Git version resolves by result (pinned: commit given, cached, remote, error)",
}, []string{"result"})

var WasmGitFetches = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_git_fetch_total",
	Help: "Commits needed for wasm pack builds by result (cached: already in the local repository, remote: fetched, error)",
}, []string{"result"})

var WasmGitCacheBytes = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_wasm_git_cache_bytes",
	Help: "Size of the local git repository cache (updated after each pack build)",
})

var WasmGitEvictions = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_wasm_git_cache_evictions_total",
	Help: "Number of repositories removed from the local git cache",
})

var WasmPackCacheBytes = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_wasm_pack_cache_bytes",
	Help: "Size of the wasm pack disk cache (updated after each build)",
})

var WasmPackEvictions = promauto.NewCounter(prometheus.CounterOpts{
	Name: "bx_wasm_pack_evictions_total",
	Help: "Number of packs removed from the wasm pack cache",
})

// wasm proxy

var WasmSessionsActive = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "bx_wasm_proxy_connections_active",
	Help: "Number of open wasm proxy websocket connections",
})

var WasmProxyConnections = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_proxy_connections_total",
	Help: "Wasm proxy connections by type (dns, http, newaddr, invalid)",
}, []string{"type"})

var WasmProxyRejected = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_proxy_rejected_total",
	Help: "Refused wasm proxy targets by stage (target, dns, http)",
}, []string{"stage"})

var WasmProxyBytes = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bx_wasm_proxy_bytes_total",
	Help: "Bytes tunneled through the wasm http proxy (in: from the client, out: to the client)",
}, []string{"direction"})
