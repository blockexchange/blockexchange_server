package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var handleHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "bx_controller_handle_hist",
	Help:    "Histogram for the controller serve time",
	Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
})

func init() {
	prometheus.Register(handleHistogram)
}
