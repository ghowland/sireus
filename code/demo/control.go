package demo

import (
	"github.com/ghowland/sireus/code/data"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
)

// If AppConfig.EnableDemo is true, this will be run in the background forever producing Prometheus data to server demonstration and educational purposes
func RunDemoForever(webPrimary *fiber.App) {
	// Start the Demo API HTTP listener
	go RunDemoAPIServer()

	// Add our API paths to the Web Primary, so we can control the demo inside the normal web app
	ConfigureDemoWebPrimary(webPrimary)

	// Set up the random number generator with a new seed
	rand.Seed(time.Now().UnixNano())

	// Prep testing time per run, so our results are time-based
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

// Configure the Demo's primary web server
func ConfigureDemoWebPrimary(webPrimary *fiber.App) {
	// Add our API paths to the Web Primary, so we can control the demo inside the normal web app
	webPrimary.Post("/demo/edge/break/circuit1", func(c *fiber.Ctx) error {
		return c.SendString("{\"_success\": " + BreakCircuit1())
	})

	webPrimary.Post("/demo/edge/break/circuit2", func(c *fiber.Ctx) error {
		return c.SendString("{\"_success\": " + BreakCircuit2())
	})

	webPrimary.Post("/demo/database/break/storage_degraded", func(c *fiber.Ctx) error {
		return c.SendString("{\"_success\": " + BreakStorageDegraded())
	})
}

// Takes the existing Render map, and adds the Demo variables to it, so they can be rendered
func UpdateRenderMapWithDemoData(renderMap map[string]interface{}) {
	renderMap["demo_current_request_per_second"] = CurrentRequestsPerSecond
}
