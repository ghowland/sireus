package demo

import "github.com/ghowland/sireus/code/data"

// If AppConfig.EnableDemo is true, this will be run in the background forever producing Prometheus data to server demonstration and educational purposes
func RunDemoForever() {
	// Run until we are quitting
	for !data.SireusData.IsQuitting {

	}
}
