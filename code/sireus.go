package main

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/data"
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

	// Run the server in the background until we end the server
	go server.RunForever()

	engine := webapp.CreateHandlebarsEngine(data.SireusData.AppConfig)
	app := webapp.CreateWebApp(engine)

	app.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(appdata.GetAPIPlotData(data.SireusData.AppConfig, c))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, data.SireusData.Site)
		return c.Render("overwatch", pageDataMap, "layouts/main_common")
	})

	app.Get("/site", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, data.SireusData.Site)
		return c.Render("site", pageDataMap, "layouts/main_common")
	})

	app.Get("/bot_group", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, data.SireusData.Site)
		return c.Render("bot_group", pageDataMap, "layouts/main_common")
	})

	app.Get("/bot", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, data.SireusData.Site)
		return c.Render("bot", pageDataMap, "layouts/main_common")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		pageDataMap := webapp.GetPageMapData(c, data.SireusData.Site)
		return c.Render("test", pageDataMap, "layouts/main_common")
	})

	// Static Files: JS, Images
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./static_web"),
	}))

	_ = app.Listen(":3000")
}
