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
var app_config_path = "config/config.json"

func main() {

	start_time := time.Now().Add(time.Duration(-60))
	prom_data := extdata.QueryPrometheus("localhost", 9090, "query_range?query=windows_service_status", start_time, 60)
	extdata.ExtractBotsFromPromData(prom_data, "name")

	app_config := appdata.LoadConfig(app_config_path)

	actionData, err := os.ReadFile((app_config.ActionPath))
	util.Check(err)

	var bot appdata.Bot
	err = json.Unmarshal([]byte(actionData), &bot)
	util.Check(err)

	engine := webapp.CreateHandlebarsEngine(app_config)

	app := webapp.CreateWebApp(engine)

	page_data_map := fiber.Map{
		"info":     "Testing 123!",
		"bot":      bot,
		"title":    "Sireus",
		"test_one": true,
		"test_two": false,
	}

	app.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(appdata.GetAPIPlotData(app_config, c))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", page_data_map, "layouts/main_common")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", page_data_map, "layouts/main_common")
	})

	// Static Files: JS, Images
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./web_static"),
	}))

	_ = app.Listen(":3000")
}
