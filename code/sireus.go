package main

import (
	"fmt"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/demo"
	"github.com/ghowland/sireus/code/exporter"
	"github.com/ghowland/sireus/code/server"
	"github.com/ghowland/sireus/code/webapp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting Sireus server...")

	// Configure the Global Server Singleton, where all the information will go
	server.Configure()

	// Run the Prometheus Exporter listener in the background
	go exporter.RunExporterListener()

	// Run the server in the background until we end the server
	go server.RunForever()

	engine := webapp.CreateHandlebarsEngine(data.SireusData.AppConfig)
	web := webapp.CreateWebApp(engine)

	// If we want to run the demo, run in the background
	if data.SireusData.AppConfig.EnableDemo {
		go demo.RunDemoForever(web)

		// Demo control through the Website.  The Demo API is for the Sireus Client, which won't have CORS issues like the Web API would
		web.Post("/demo/action", func(c *fiber.Ctx) error {
			return c.SendString(demo.ProcessWebDemoAction(c))
		})
	}

	// API Calls
	web.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(app.GetAPIPlotData(c))
	})

	web.Post("/api/plot_metrics", func(c *fiber.Ctx) error {
		return c.SendString(app.GetAPIPlotMetrics(c))
	})

	web.Post("/api/web/bot", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromRPC(c, &data.SireusData.Site)
		return c.SendString(webapp.RenderRPCHtml("web/bot.hbs", renderMap))
	})

	web.Post("/api/web/demo_control", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromRPC(c, &data.SireusData.Site)
		return c.SendString(webapp.RenderRPCHtml("web/demo_control.hbs", renderMap))
	})

	// Raw Data Page - Not an API call
	web.Get("/raw/metrics", func(c *fiber.Ctx) error {
		return c.SendString(app.GetRawMetricsJSON(c))
	})

	// Web Pages
	web.Get("/", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("demo_control", renderMap, "layouts/main_common")
	})

	web.Get("/site", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("site", renderMap, "layouts/main_common")
	})

	web.Get("/bot_group", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("bot_group", renderMap, "layouts/main_common")
	})

	web.Get("/bot", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("bot", renderMap, "layouts/main_common")
	})

	web.Get("/site_query", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("site_query", renderMap, "layouts/main_common")
	})

	web.Get("/overwatch", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("overwatch", renderMap, "layouts/main_common")
	})

	web.Get("/show_prom", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("show_prom", renderMap, "layouts/main_common")
	})

	web.Get("/show_config", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("show_config", renderMap, "layouts/main_common")
	})

	web.Get("/demo_info", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("demo_info", renderMap, "layouts/main_common")
	})

	web.Get("/test", func(c *fiber.Ctx) error {
		renderMap := webapp.GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("test", renderMap, "layouts/main_common")
	})

	// Static Files: JS, Images
	web.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./static_web"),
	}))

	_ = web.Listen(fmt.Sprintf(":%d", data.SireusData.AppConfig.WebHttpPort))
}
