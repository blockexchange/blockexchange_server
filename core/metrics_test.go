package core_test

import (
	"blockexchange/core"
	"blockexchange/metrics"
	"blockexchange/types"
	"context"
	"path/filepath"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// returns a function that gives the increase of the counter since this call
func delta(c prometheus.Collector) func() float64 {
	before := testutil.ToFloat64(c)
	return func() float64 { return testutil.ToFloat64(c) - before }
}

func TestWasmPackMetrics(t *testing.T) {
	c := newPackCore(t)
	game_dir := filepath.Join(t.TempDir(), "game")
	game_commit := commitFiles(t, game_dir, map[string]string{"game.conf": "v1\n"})

	kind := string(core.WasmPackKindDependencies)
	miss := delta(metrics.WasmPackRequests.WithLabelValues(kind, "miss"))
	hit := delta(metrics.WasmPackRequests.WithLabelValues(kind, "hit"))
	pinned := delta(metrics.WasmPackResolves.WithLabelValues("pinned"))
	remote := delta(metrics.WasmPackResolves.WithLabelValues("remote"))
	cached := delta(metrics.WasmPackResolves.WithLabelValues("cached"))

	latest := []*types.SchemaModDependency{dep(types.SchemaModDependencyTypeGame, game_dir, "latest")}
	_, err := c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, latest)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, miss())
	assert.Equal(t, 1.0, remote())
	assert.Greater(t, testutil.ToFloat64(metrics.WasmPackCacheBytes), 0.0)

	_, err = c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, latest)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, hit())
	assert.Equal(t, 1.0, cached())

	_, err = c.GetWasmPack(context.Background(), core.WasmPackKindDependencies, []*types.SchemaModDependency{
		dep(types.SchemaModDependencyTypeGame, game_dir, game_commit),
	})
	assert.NoError(t, err)
	assert.Equal(t, 2.0, hit())
	assert.Equal(t, 1.0, pinned())
	assert.Equal(t, 1.0, miss())
}

func TestWasmProxyMetrics(t *testing.T) {
	url, _ := setupWasmProxy(t)

	http_conns := delta(metrics.WasmProxyConnections.WithLabelValues("http"))
	bytes_in := delta(metrics.WasmProxyBytes.WithLabelValues("in"))
	bytes_out := delta(metrics.WasmProxyBytes.WithLabelValues("out"))
	rejected_target := delta(metrics.WasmProxyRejected.WithLabelValues("target"))
	rejected_http := delta(metrics.WasmProxyRejected.WithLabelValues("http"))

	ws := dialProxy(t, url, "PROXY IPV4 TCP 10.0.0.1 8080")
	expectOK(t, ws)
	assert.NoError(t, ws.WriteMessage(websocket.BinaryMessage, []byte("GET http://example.com/ HTTP/1.1\r\nHost: example.com\r\n\r\n")))
	assert.Equal(t, 403, readResponse(t, ws).StatusCode)
	assert.Equal(t, 1.0, http_conns())
	assert.Equal(t, 1.0, rejected_http())
	assert.Greater(t, bytes_in(), 0.0)
	assert.Greater(t, bytes_out(), 0.0)

	ws = dialProxy(t, url, "PROXY IPV4 UDP 1.2.3.4 30000")
	_, _, err := ws.ReadMessage()
	assert.Error(t, err)
	assert.Equal(t, 1.0, rejected_target())

	// all connections closed
	assert.Eventually(t, func() bool { return testutil.ToFloat64(metrics.WasmSessionsActive) == 0 }, 5e9, 1e7)
}
