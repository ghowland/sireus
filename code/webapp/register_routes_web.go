package webapp

import (
	"fmt"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/demo"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
)

// Register the routes for the web pages
func RegisterRoutesWeb(web *fiber.App) {

	// Raw Data Page - Not an API call
	web.Get("/raw/metrics", func(c *fiber.Ctx) error {
		return c.SendString(app.GetRawMetricsJSON(c))
	})

	// Web Pages
	web.Get("/", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		// Update the Demo Control with demo specific data
		demo.UpdateRenderMapWithDemoData(renderMap)
		return c.Render("demo_control", renderMap, "layouts/main_common")
	})

	web.Get("/site", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("site", renderMap, "layouts/main_common")
	})

	web.Get("/bot_group", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("bot_group", renderMap, "layouts/main_common")
	})

	web.Get("/bot", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("bot", renderMap, "layouts/main_common")
	})

	web.Get("/site_query", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("site_query", renderMap, "layouts/main_common")
	})

	web.Get("/overwatch", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("overwatch", renderMap, "layouts/main_common")
	})

	web.Get("/show_prom", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		url := fmt.Sprintf("http://localhost:%d/metrics", data.SireusData.AppConfig.PrometheusExportPort)
		body, err := util.HttpGet(url)
		if util.Check(err) {
			body = fmt.Sprintf("Failed to get Prometheus data from URL: %s  Error: %s", url, err.Error())
		}
		renderMap["prometheus_exporter"] = body
		return c.Render("show_prom", renderMap, "layouts/main_common")
	})

	web.Get("/show_config", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("show_config", renderMap, "layouts/main_common")
	})

	web.Get("/demo_info", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("demo_info", renderMap, "layouts/main_common")
	})

	web.Get("/test", func(c *fiber.Ctx) error {
		renderMap := GetRenderMapFromParams(c, &data.SireusData.Site)
		return c.Render("test", renderMap, "layouts/main_common")
	})
}
