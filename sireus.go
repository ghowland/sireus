package main

import (
	"encoding/json"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	"os"
)

type ActionConsideration struct {
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	CurveName  string  `json:"curve"`
	RangeStart float32 `json:"range_start"`
	RangeEnd   float32 `json:"range_end"`
}

type Action struct {
	Name           string                `json:"name"`
	Weight         float32               `json:"weight"`
	WeightMin      float32               `json:"weight_min"`
	Considerations []ActionConsideration `json:"considerations"`
}

type Bot struct {
	Name    string   `json:"name"`
	Actions []Action `json:"actions"`
}

type CurveData struct {
	Name   string    `json:"name"`
	Values []float32 `json:"values"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetActionsHtml(bot Bot) string {
	output := ""

	output += fmt.Sprintf("Bot: %s\n", bot.Name)

	for _, action := range bot.Actions {
		output += fmt.Sprintf("Action: %s  Weight: %.2f Min: %.2f\n", action.Name, action.Weight, action.WeightMin)

		for _, consider := range action.Considerations {
			output += fmt.Sprintf("  Consider: %s  Weight: %.2f  Curve: %s  Range Start: %.2f  End: %.2f\n", consider.Name, consider.Weight, consider.CurveName, consider.RangeStart, consider.RangeEnd)
		}

	}

	return output
}

type AppConfig struct {
	WebPath         string `json:"web_path"`
	ActionPath      string `json:"action_path"`
	CurvePathFormat string `json:"curve_path_format"`
}

var app_config_path = "config/config.json"

func LoadConfig() AppConfig {
	app_config_data, err := os.ReadFile((app_config_path))
	Check(err)

	var app_config AppConfig
	json.Unmarshal([]byte(app_config_data), &app_config)

	return app_config
}

func LoadCurveData(path string) CurveData {
	curveData, err := os.ReadFile((path))
	Check(err)

	var curve_data CurveData
	json.Unmarshal([]byte(curveData), &curve_data)

	return curve_data
}

func CreateHandlebarsEngine(app_config AppConfig) *handlebars.Engine {
	// Handlebars Engine for Fiber
	engine := handlebars.New(app_config.WebPath, ".hbs")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	//// Debug will print each template that is parsed, good for debugging
	//engine.Debug(true) // Optional. Default: false

	//// Layout defines the variable name that is used to yield templates within layouts
	//engine.Layout("embed") // Optional. Default: "embed"

	raymond.RegisterHelper("botinfo", func(bot Bot) string {
		return bot.Name + "  Actions: " + string(len(bot.Actions))
	})

	return engine
}

func CreateWebApp(engine *handlebars.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return app
}

func main() {
	app_config := LoadConfig()

	actionData, err := os.ReadFile((app_config.ActionPath))
	Check(err)

	var bot Bot
	json.Unmarshal([]byte(actionData), &bot)

	curve_path := fmt.Sprintf(app_config.CurvePathFormat, bot.Actions[0].Considerations[0].CurveName)
	curve_data := LoadCurveData(curve_path)

	engine := CreateHandlebarsEngine(app_config)

	app := CreateWebApp(engine)

	app.Get("/", func(c *fiber.Ctx) error {
		//return c.SendString(GetActionsHtml(bot))
		return c.Render("index", fiber.Map{
			"info":       "Testing 123!",
			"bot":        bot,
			"title":      "Sireus",
			"curve_data": curve_data,
		}, "layouts/main_common")
	})

	app.Listen(":3000")

}
