package prompusher

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NewGauge https://prometheus.io/docs/concepts/metric_types/#gauge
// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
func (c *client) NewGauge(name, help string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: "etzba",
		Name:      name,
		Help:      help,
	}, labels)
	return gauge
}

// NewCounter https://prometheus.io/docs/concepts/metric_types/#counter
// A counter is a cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart.
// Use to count total executions from a result
func (c *client) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: c.Namespace,
		Subsystem: "etzba",
		Name:      name,
		Help:      help,
	}, labels)
	return counter
}

// NewHistogram https://prometheus.io/docs/concepts/metric_types/#histogram
// A histogram samples observations (usually things like request durations or response sizes) and counts them in configurable buckets. It also provides a sum of all observed values.
// The best fit for this tool to collect results.
func (c *client) NewHistogram(name, help string, labels []string) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: c.Namespace,
		Subsystem: "etzba",
		Name:      name,
		Help:      help,
	}, labels)
	return histogram
}
