package demo

import (
	"fmt"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// This is the API server for the Sireus Client to use.  The web app demo controls do not use this, they go through the web app to avoid CORS issues.
func RunDemoAPIServer() {
	app := fiber.New(fiber.Config{})

	// Fixing issues in the demo
	app.Get("/fix/circuit", func(c *fiber.Ctx) error {
		name := c.Query("name", "")
		//log.Printf("Demo API: Fix Circuit: '%s'", name)
		switch name {
		case "SFO-WAS-11":
			FixCircuit2()
			return c.SendString("{\"status\": 201}")
		case "SFO-LAS-27":
			//FixCircuit1()		// Must be fixed manually by the user for the demo
			return c.SendString("{\"status\": 501}") // Failure
		default:
			return c.SendString("{\"status\": 404}") // Failure
		}
	})

	app.Get("/fix/database_storage_degraded", func(c *fiber.Ctx) error {
		success := FixStorageDegraded()
		if success {
			return c.SendString("{\"status\": 200}")
		} else {
			return c.SendString("{\"status\": 505}")
		}
	})

	_ = app.Listen(fmt.Sprintf(":%d", data.SireusData.AppConfig.DemoApiPort))
}

// Process RPC requests from the web app to wrap Demo functionality.  This is a different path than the Sireus Client uses with RunDemoAPIServer()
func ProcessWebDemoAction(c *fiber.Ctx) string {
	input := util.ParseContextBody(c)

	actionName, ok := input["action"]
	if !ok {
		return "{\"_failure\": \"Invalid Demo Action request, no 'action' field specified.\"}"
	}

	botGroupName, ok := input["bot_group"]
	botName, ok := input["bot"]

	switch actionName {
	case "break":
		return DemoBreakBot(botGroupName, botName)
	case "fix":
		return DemoFixBot(botGroupName, botName)
	case "clear_command_history":
		app.AdminClearCommandHistory()
		return fmt.Sprintf("{\"_success\": \"Command History has been cleared.\"}")
	case "set_edge_traffic":
		valueStr, ok := input["value"]
		if !ok {
			return fmt.Sprintf("{\"_failure\": \"Missing key 'value'\"}")
		}
		value, err := strconv.ParseFloat(valueStr, 64)
		if util.Check(err) {
			return fmt.Sprintf("{\"_failure\": \"Invalid value: %s\"}", err.Error())
		}
		// Set the demo requests per second
		CurrentRequestsPerSecond = value
		return fmt.Sprintf("{\"_success\": \"Edge traffic set to %0.0f requests per second\"}", CurrentRequestsPerSecond)
	default:
		return fmt.Sprintf("{\"_failure\": \"Unknown action: %s\"}", actionName)
	}
}

// Create problems in the Demo, which will cause the metrics to be updated and Sireus will respond
func DemoBreakBot(botGroupName string, botName string) string {
	message := ""

	switch botGroupName {
	case "Edge":
		switch botName {
		case "SFO-WAS-11":
			BreakCircuit2()
			message = "Breaking circuit SFO-WAS-11, please wait about 15 seconds and check Action History"
			break
		case "SFO-LAS-27":
			BreakCircuit1()
			message = "Breaking circuit SFO-LAS-27, please wait about 15 seconds and check Action History"
			break
		}
		break
	case "Database":
		BreakStorageDegraded()
		message = "Degrading database storage, please wait about 15 seconds and check Action History"
		break
	}

	return fmt.Sprintf("{\"_success\": \"%s\"}", message)
}

// Fix the problems in the Demo, which will cause the metrics to be updated and Sireus will respond
func DemoFixBot(botGroupName string, botName string) string {
	message := ""

	switch botGroupName {
	case "Edge":
		switch botName {
		case "SFO-WAS-11":
			FixCircuit2()
			message = "Fixing circuit SFO-WAS-11, please wait about 15 seconds and check Edge States"
			break
		case "SFO-LAS-27":
			FixCircuit1()
			message = "Fixing circuit SFO-LAS-27, please wait about 15 seconds and check Edge States"
			break
		}
		break
	case "Database":
		FixStorageDegraded()
		message = "Fixing database storage, please wait about 15 seconds and check Edge States"
		break
	}

	return fmt.Sprintf("{\"_success\": \"%s\"}", message)
}
