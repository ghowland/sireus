package webapp

import (
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/demo"
	"github.com/gofiber/fiber/v2"
)

// Register the Routes for the Web RPC API
func RegisterRoutesAPI(web *fiber.App) {

	// API Calls
	web.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(app.GetAPIPlotData(c))
	})

	web.Post("/api/plot_metrics", func(c *fiber.Ctx) error {
		return c.SendString(app.GetAPIPlotMetrics(c))
	})

	web.Post("/api/web/bot", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromRPC(c, &data.SireusData.Site)
		return c.SendString(RenderRPCHtml("web/bot.hbs", renderMap))
	})

	web.Post("/api/web/demo_control", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromRPC(c, &data.SireusData.Site)
		// Update the Demo Control with demo specific data
		demo.UpdateRenderMapWithDemoData(renderMap)
		return c.SendString(RenderRPCHtml("web/demo_control.hbs", renderMap))
	})

}
