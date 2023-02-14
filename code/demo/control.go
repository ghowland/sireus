package demo

import (
	"github.com/ghowland/sireus/code/data"
	"math/rand"
	"time"
)

// If AppConfig.EnableDemo is true, this will be run in the background forever producing Prometheus data to server demonstration and educational purposes
func RunDemoForever() {
	rand.Seed(time.Now().UnixNano())

	lastRunTime := time.Now()

	// Run until we are quitting
	for !data.SireusData.IsQuitting {
		// Get time since last run, so all the flow is time based
		currentRunTime := time.Now()
		secondsSinceLastRun := currentRunTime.Sub(lastRunTime).Seconds()

		// Update all the components variables
		UpdateEdge(secondsSinceLastRun)
		UpdateApp(secondsSinceLastRun)
		UpdateDatabase(secondsSinceLastRun)

		// Give back time
		time.Sleep(200 * time.Millisecond)
		lastRunTime = currentRunTime
	}
}
