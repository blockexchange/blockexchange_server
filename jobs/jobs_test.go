package jobs

import (
	"blockexchange/metrics"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRunJobMetrics(t *testing.T) {
	ok := metrics.JobRuns.WithLabelValues("test-job", "ok")
	failed := metrics.JobRuns.WithLabelValues("test-job", "error")
	last := metrics.JobLastSuccess.WithLabelValues("test-job")

	runJob("test-job", func() error { return errors.New("broken") })
	assert.Equal(t, 1.0, testutil.ToFloat64(failed))
	assert.Equal(t, 0.0, testutil.ToFloat64(ok))
	assert.Equal(t, 0.0, testutil.ToFloat64(last))

	runJob("test-job", func() error { return nil })
	assert.Equal(t, 1.0, testutil.ToFloat64(ok))
	assert.Greater(t, testutil.ToFloat64(last), 0.0)
}
