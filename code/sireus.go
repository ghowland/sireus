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
	log.Printf("OUTSIDE: Bots after Prom Update: %s  Count: %d", site.BotGroups[0].Name, len(site.BotGroups[0].Bots))

	engine := webapp.CreateHandlebarsEngine(appConfig)

	app := webapp.CreateWebApp(engine)

	pageDataMap := fiber.Map{
		"site":     site,
		"botGroup": site.BotGroups[0],
		"title":    "Sireus",
	}

	app.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(appdata.GetAPIPlotData(appConfig, c))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", pageDataMap, "layouts/main_common")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", pageDataMap, "layouts/main_common")
	})

	// Static Files: JS, Images
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./web_static"),
	}))

	_ = app.Listen(":3000")
}
