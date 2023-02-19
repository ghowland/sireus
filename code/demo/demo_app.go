package demo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	// What is the state of the Demo App?
	AppState int64
)

const (
	AppNormal AppState = iota // Demo App is normal
)

var (
	// App simulation: Wait Queue
	appWaiting = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_req_queue_wait",
		Help: "The current requests waiting to be processed",
		ConstLabels: map[string]string{
			"bot_group": "App",
		},
	})

	// App simulation: Timeouts
	appTimeout = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_req_timeout",
		Help: "The total requests that have timed out",
		ConstLabels: map[string]string{
			"bot_group": "App",
		},
	})

	// App simulation: Success
	appSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_req_success",
		Help: "The total requests that have been processed successfully",
		ConstLabels: map[string]string{
			"bot_group": "App",
		},
	})
)

var (
	AppRequestQueueLength int = 0
)

// Update the Demo App
func UpdateApp(seconds float64) {
	// This demo app is stateless, as it should be
}

// Receive request timeouts from the database.  This might feel backwards, but its a demo simulation
func ReceiveTimeoutsFromDatabase(requests int) {
	appTimeout.Add(float64(requests))

	// Update the wait queue
	AppRequestQueueLength -= requests
	if AppRequestQueueLength < 0 {
		AppRequestQueueLength = 0
	}
	appWaiting.Set(float64(AppRequestQueueLength))
}

// Receive request successes from the database.  This might feel backwards, but its a demo simulation
func ReceiveSuccessFromDatabase(requests int) {
	appSuccess.Add(float64(requests))

	// Send our successful requests back to the edge to be delivered to requester
	ReceiveSuccessFromApp(requests)

	// Update the wait queue
	AppRequestQueueLength -= requests
	if AppRequestQueueLength < 0 {
		AppRequestQueueLength = 0
	}
	appWaiting.Set(float64(AppRequestQueueLength))
}

// Receive requests from the edge.
func ReceiveRequestsFromEdge(requests int) {
	AppRequestQueueLength += requests
	if AppRequestQueueLength < 0 {
		AppRequestQueueLength = 0
	}
	appWaiting.Set(float64(AppRequestQueueLength))

	// Send the requests to the database
	AddDatabaseRequests(requests)
}
