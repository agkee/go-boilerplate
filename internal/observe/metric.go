package observe

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	SomeMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "some_metric",
			Help: "A counter for showing how to setup prometheus",
		},
		[]string{},
	)
)

// RegisterPrometheus adds the prometheus handler to the mux router
// Notethat every metric has to be registered with prometheus for it show up
// when the /metrics route is hit.
func RegisterPrometheus(m *mux.Router) *mux.Router {
	prometheus.MustRegister(SomeMetric)

	m.Handle("/metrics", promhttp.Handler())
	return m
}
