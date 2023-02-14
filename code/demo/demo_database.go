package demo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	// What is the state of the Demo Database?
	DatabaseState int64
)

const (
	DatabaseNormal          DatabaseState = iota // Demo Database is normal
	DatabaseStorageDegraded                      // Demo Database has degraded storage
)

var (
	// Database simulation: Wait Queue
	databaseWaiting = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "demo_database_req_queue_wait",
		Help:        "The current requests waiting to be processed",
		ConstLabels: map[string]string{},
	})

	// Database simulation: Timeouts
	databaseTimeout = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "demo_database_req_timeout",
		Help:        "The total requests that have timed out",
		ConstLabels: map[string]string{},
	})

	// Database simulation: Requests Successful
	databaseSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "demo_database_req_success",
		Help:        "The total requests that have been processed successfully",
		ConstLabels: map[string]string{},
	})
)

var (
	// Simulating request processing: Normal Speed
	NormalRequestProcessSpeed int = 1000

	// Simulating request processing: Degraded Speed
	DegradedRequestProcessSpeed int = 250

	// Simulating request processing: Timeout over N requests in the queue, anything over this will be timed out to balance the system
	TimeoutWaitLimit int = 2000

	// Current Queue Length, this gets exported to Prometheus as databaseWaiting
	DatabaseRequestQueueLength int = 0

	// Current Database State
	CurrentDatabaseState DatabaseState
)

func PerSecond(original int, seconds float64) int {
	value := float64(original) * seconds

	return int(value)
}

// Update the Demo Database
func UpdateDatabase(seconds float64) {
	// Simulate Processing requests, at a normal and storage-degraded speed
	processSpeed := PerSecond(NormalRequestProcessSpeed, seconds)
	if CurrentDatabaseState == DatabaseStorageDegraded {
		processSpeed = PerSecond(DegradedRequestProcessSpeed, seconds)
	}

	// If all queued requests are less than our processing speed, then we process them all
	if DatabaseRequestQueueLength < processSpeed {
		// Send the Demo App our success count
		ReceiveSuccessFromDatabase(DatabaseRequestQueueLength)

		// All are processed
		databaseSuccess.Add(float64(DatabaseRequestQueueLength))
		// Nothing in the queue
		DatabaseRequestQueueLength = 0
	} else {
		// Send the Demo App our success count
		ReceiveSuccessFromDatabase(processSpeed)

		// Else, process our current rate as successful
		databaseSuccess.Add(float64(processSpeed))

		// Subtract that rate from the queue.  It will remain positive
		DatabaseRequestQueueLength -= processSpeed
	}

	// Deal with Timeouts
	if DatabaseRequestQueueLength > TimeoutWaitLimit {
		// Determine the number of timeouts.  This is not realistic, but it simulates and balances the system
		timeoutCount := DatabaseRequestQueueLength - TimeoutWaitLimit

		// Timeouts per second, will slow down how many we do immediately, for some of them to be handled
		//NOTE(ghowland): This is not a good simulation, but it is way easier to implement, and simulates a simulation
		timeoutNow := PerSecond(timeoutCount, seconds)

		// Upper bound timeouts, because math
		if timeoutNow > DatabaseRequestQueueLength {
			timeoutNow = DatabaseRequestQueueLength
		}

		// Add these to our database timeouts
		databaseTimeout.Add(float64(timeoutNow))

		// Send timeouts back to the Demo App server
		ReceiveTimeoutsFromDatabase(timeoutNow)

		// Remove the timeout count from our request queue
		DatabaseRequestQueueLength -= timeoutNow
	}

	// Update the request queue
	databaseWaiting.Set(float64(DatabaseRequestQueueLength))
}

// Receive demo requests from the App server
func AddDatabaseRequests(requests int) {
	DatabaseRequestQueueLength += requests
}
