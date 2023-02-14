package demo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"math/rand"
)

type (
	// What is the state of the Demo Edge?
	EdgeState int64
)

const (
	EdgeNormal         EdgeState = iota // Demo Edge is normal
	EdgeDownCircuit1                    // Demo Edge Circuit 1 is down
	EdgeDownCircuit2                    // Demo Edge Circuit 2 is down
	EdgeDownCircuitAll                  // Demo Edge All Circuits are down
)

var (
	edgeRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_requests",
		Help: "The total number of requests sent on the Internet to this service.  This is demo info, not simulation.",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	// Circuit 1: SFO-LAS-27
	edgeDataIn1 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_in_octets",
		Help: "The total number of bytes received",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	edgeDataOut1 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_out_octets",
		Help: "The total number of bytes sent",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	edgeLinkState1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_edge_if_link_state",
		Help: "The link state of this connection",
		ConstLabels: map[string]string{
			"circuit": "SFO-LAS-27",
		},
	})

	// Circuit 2: SFO-LAS-27
	edgeDataIn2 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_in_octets",
		Help: "The total number of bytes received",
		ConstLabels: map[string]string{
			"circuit": "SFO-WAS-11",
		},
	})

	edgeDataOut2 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_edge_if_out_octets",
		Help: "The total number of bytes sent",
		ConstLabels: map[string]string{
			"circuit": "SFO-WAS-11",
		},
	})

	edgeLinkState2 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_edge_if_link_state",
		Help: "The link state of this connection",
		ConstLabels: map[string]string{
			"circuit": "SFO-WAS-11",
		},
	})
)

var (
	CurrentEdgeState EdgeState

	NormalRequestBaseIn  int = 100
	NormalRequestBaseOut int = 100 * 7

	NormalRequestRandomIn  int = 250
	NormalRequestRandomOut int = 250 * 5

	// Demo Control: How many requests per second are coming into the system?
	CurrentRequestsPerSecond float64 = 800
)

// Update the Demo Edge
func UpdateEdge(seconds float64) {
	// We get new demo requests every N seconds (0.2), and then update all our counters based on randomized traffic size
	requestAdded := int(CurrentRequestsPerSecond * seconds)
	edgeRequests.Add(float64(requestAdded))

	// Process the ingress only.  Have to wait for successes to process egress, for the purpose of the demo.
	//NOTE(ghowland): Not sending ACK type data back on purpose here, better it's a clean 0 for demo purposes.
	//		Also, not simulating things like bandwidth caps, etc.  Most basic simulation of pass or no-pass.
	switch CurrentEdgeState {
	case EdgeNormal:
		// Perform a random range calculate Base + Random, then multiply it by number of requests we are adding.  Split by 2 circuits.
		edgeDataIn1.Add(float64(rand.Intn(NormalRequestRandomIn)+NormalRequestBaseIn) * float64(requestAdded/2))
		edgeLinkState1.Set(1)

		edgeDataIn2.Add(float64(rand.Intn(NormalRequestRandomIn)+NormalRequestBaseIn) * float64(requestAdded/2))
		edgeLinkState2.Set(1)

		// Add the requests to the Demo App
		ReceiveRequestsFromEdge(requestAdded)
		break
	case EdgeDownCircuit1:
		edgeLinkState1.Set(0)

		edgeDataIn2.Add(float64(rand.Intn(NormalRequestRandomIn*2)+NormalRequestBaseIn*2) * float64(requestAdded))
		edgeLinkState2.Set(1)

		// Add the requests to the Demo App
		ReceiveRequestsFromEdge(requestAdded)
		break
	case EdgeDownCircuit2:
		edgeDataIn1.Add(float64(rand.Intn(NormalRequestRandomIn*2)+NormalRequestBaseIn*2) * float64(requestAdded))
		edgeLinkState1.Set(1)

		edgeLinkState2.Set(0)

		// Add the requests to the Demo App
		ReceiveRequestsFromEdge(requestAdded)
		break
	case EdgeDownCircuitAll:
		edgeLinkState1.Set(0)
		edgeLinkState2.Set(0)

		// No requests to add to the Demo App, because both circuits are down
		break
	}
}

func ReceiveSuccessFromApp(requests int) {
	// Process the requests through their state
	switch CurrentEdgeState {
	case EdgeNormal:
		// Perform a random range calculate Base + Random, then multiply it by number of requests we are adding.  Split by 2 circuits.
		edgeDataOut1.Add(float64(rand.Intn(NormalRequestRandomOut)+NormalRequestBaseOut) * float64(requests/2))
		edgeDataOut2.Add(float64(rand.Intn(NormalRequestRandomOut)+NormalRequestBaseOut) * float64(requests/2))
		break
	case EdgeDownCircuit1:
		edgeDataOut2.Add(float64(rand.Intn(NormalRequestRandomOut*2)+NormalRequestBaseOut*2) * float64(requests))
		break
	case EdgeDownCircuit2:
		edgeDataOut1.Add(float64(rand.Intn(NormalRequestRandomOut*2)+NormalRequestBaseOut*2) * float64(requests))
		break
	case EdgeDownCircuitAll:
		// No requests can be sent back.  All circuits are down
		break
	}
}

// Dont delay, as would be normal, just immediately bring the circuit up to make the interactive demo move faster
func FixCircuit1() {
	if CurrentEdgeState == EdgeDownCircuit1 {
		CurrentEdgeState = EdgeNormal
	} else if CurrentEdgeState == EdgeDownCircuitAll {
		CurrentEdgeState = EdgeDownCircuit2
	}
}

// Dont delay, as would be normal, just immediately bring the circuit up to make the interactive demo move faster
func FixCircuit2() {
	if CurrentEdgeState == EdgeDownCircuit2 {
		CurrentEdgeState = EdgeNormal
	} else if CurrentEdgeState == EdgeDownCircuitAll {
		CurrentEdgeState = EdgeDownCircuit1
	}
}
