package main

import (
	"encoding/json"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/extdata"
	"github.com/ghowland/sireus/code/util"
	"github.com/ghowland/sireus/code/webapp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	"os"
	"time"
)

// TODO(ghowland): Handle config file as CLI arg
var appConfigPath = "config/config.json"

func main() {

	startTime := time.Now().Add(time.Duration(-60))
	promData := extdata.QueryPrometheus("localhost", 9090, "query_range?query=windows_service_status", startTime, 60)
	extdata.ExtractBotsFromPromData(promData, "name")

	appConfig := appdata.LoadConfig(appConfigPath)

	actionData, err := os.ReadFile(appConfig.ActionPath)
	util.Check(err)

	var bot appdata.Bot
	err = json.Unmarshal(actionData, &bot)
	util.Check(err)

	engine := webapp.CreateHandlebarsEngine(appConfig)

	app := webapp.CreateWebApp(engine)

	pageDataMap := fiber.Map{
		"bot":      bot,
		"title":    "Sireus",
		"test_one": true,
		"test_two": false,
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
