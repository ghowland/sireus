package data

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Dynamically creates Guage metric exporters
type PrometheusMetricGauge struct {
	Key    string           // String to look up this metric
	Metric prometheus.Gauge // This is what we export with, created with ConstLabels at start, which cant change
}

// Dynamically creates Counter metric exporters
type PrometheusMetricCounter struct {
	Key    string             // String to look up this metric
	Metric prometheus.Counter // This is what we export with, created with ConstLabels at start, which cant change
}

// Managed the dynamically created Gauges and Counters
type PrometheusExportData struct {
	Gauges   []PrometheusMetricGauge   // Track all the gauges for our bots, use for variables or things that we want to save state on
	Counters []PrometheusMetricCounter // Track all the gauges for our bots, use for throwaways, where we don't care about current value of the data, just tracking it for some later process that does care
}
