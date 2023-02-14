package demo

import (
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/gofiber/fiber/v2"
)

func RunDemoAPIServer() {
	app := fiber.New(fiber.Config{})

	// Fixing issues in the demo
	app.Get("/fix/circuit1", func(c *fiber.Ctx) error {
		FixCircuit1()
		return c.SendString("{\"status\": 201}")
	})

	app.Get("/fix/circuit2", func(c *fiber.Ctx) error {
		FixCircuit2()
		return c.SendString("{\"status\": 202}")
	})

	app.Get("/fix/database_storage_degraded", func(c *fiber.Ctx) error {
		FixStorageDegraded()
		return c.SendString("{\"status\": 200}")
	})

	_ = app.Listen(fmt.Sprintf(":%d", data.SireusData.AppConfig.DemoApiPort))
}
