package main

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/extdata"
	"github.com/ghowland/sireus/code/webapp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
)

// TODO(ghowland): Handle config file as CLI arg
var appConfigPath = "config/config.json"

func main() {
	log.Print("Starting Sireus server...")

	appConfig := appdata.LoadConfig(appConfigPath)

	site := appdata.LoadSiteConfig(appConfig)
	extdata.UpdateSiteBotGroups(&site)

	engine := webapp.CreateHandlebarsEngine(appConfig)
	app := webapp.CreateWebApp(engine)

	app.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(appdata.GetAPIPlotData(appConfig, c))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, site)
		return c.Render("overwatch", pageDataMap, "layouts/main_common")
	})

	app.Get("/site", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, site)
		return c.Render("site", pageDataMap, "layouts/main_common")
	})

	app.Get("/bot_group", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, site)
		return c.Render("bot_group", pageDataMap, "layouts/main_common")
	})

	app.Get("/bot", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, site)
		return c.Render("bot", pageDataMap, "layouts/main_common")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, site)
		return c.Render("test", pageDataMap, "layouts/main_common")
	})

	// Static Files: JS, Images
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./static_web"),
	}))

	_ = app.Listen(":3000")
}
