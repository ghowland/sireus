package demo

import (
	"github.com/ghowland/sireus/code/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	edgeDataIn = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_in_octets",
		Help: "The total number of bytes received",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	edgeDataOut = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_out_octets",
		Help: "The total number of bytes sent",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	edgeLinkState = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_edge_if_link_state",
		Help: "The link state of this connection",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})
)

// If AppConfig.EnableDemo is true, this will be run in the background forever producing Prometheus data to server demonstration and educational purposes
func RunDemoForever() {
	// Run until we are quitting
	for !data.SireusData.IsQuitting {
		edgeDataIn.Add(5)
		edgeDataOut.Add(13)
		edgeLinkState.Set(1)

		// Give back time
		time.Sleep(200 * time.Millisecond)
	}
}
