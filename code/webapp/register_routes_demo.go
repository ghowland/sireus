package webapp

import (
	"github.com/ghowland/sireus/code/demo"
	"github.com/gofiber/fiber/v2"
)

// Register the routes for the demo, which is enabled or not in the config.  Optional.
func RegisterRoutesDemo(web *fiber.App) {
	// Demo control through the Website.  The Demo API is for the Sireus Client, which won't have CORS issues like the Web API would
	web.Post("/demo/action", func(c *fiber.Ctx) error {
		return c.SendString(demo.ProcessWebDemoAction(c))
	})
}
