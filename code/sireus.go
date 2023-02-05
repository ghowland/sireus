package main

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/extdata"
	"github.com/ghowland/sireus/code/webapp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
	"time"
)

// TODO(ghowland): Handle config file as CLI arg
var appConfigPath = "config/config.json"

func main() {
	log.Print("Starting Sireus server...")
	
	startTime := time.Now().Add(time.Duration(-60))
	promData := extdata.QueryPrometheus("localhost", 9090, "query_range?query=windows_service_status", startTime, 60)
	extdata.ExtractBotsFromPromData(promData, "name")

	appConfig := appdata.LoadConfig(appConfigPath)

	site := appdata.LoadSiteConfig(appConfig)

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
