package main

import (
	"encoding/json"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	Info    string   `json:"info"`
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

func LoadCurveData(app_config AppConfig, name string) CurveData {
	path := fmt.Sprintf(app_config.CurvePathFormat, name)

	curveData, err := os.ReadFile((path))
	Check(err)

	var curve_data CurveData
	json.Unmarshal([]byte(curveData), &curve_data)

	//log.Println("Load Curve Data: ", curve_data.Name)

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

	raymond.RegisterHelper("ifconsiderlength", func(considerations []ActionConsideration, count int, options *raymond.Options) raymond.SafeString {
		//log.Println("ifconsiderlength: ", len(considerations), " Count: ", count)

		if len(considerations) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	return engine
}

func CreateWebApp(engine *handlebars.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return app
}

func GetCurveDataX(curve_data CurveData) []float32 {
	var x_array []float32

	for i := 0; i < len(curve_data.Values); i++ {
		x_array = append(x_array, float32(i)*0.01)
	}

	return x_array
}

func GetCurveValue(curve_data CurveData, x float32) float32 {

	for i := 0; i < len(curve_data.Values); i++ {
		cur_pos_x := float32(i) * 0.01
		if x <= cur_pos_x {
			return curve_data.Values[i]
		}
	}

	return curve_data.Values[len(curve_data.Values)-1]
}

func ParseContextBody(c *fiber.Ctx) map[string]string {
	values, err := url.ParseQuery(string(c.Body()))
	Check(err)

	obj := map[string]string{}
	for k, v := range values {
		if len(v) > 0 {
			obj[k] = v[0]
		}
	}

	return obj
}

func GetAPIPlotData(app_config AppConfig, c *fiber.Ctx) string {
	input := ParseContextBody(c)
	//log.Println("Get API Plot Data: ", input)

	if input["name"] == "" {
		failure_result := map[string]interface{}{
			"_failure": "Name not found, aborting",
		}
		failure_json, _ := json.Marshal(failure_result)
		return string(failure_json)
	}

	curve_data := LoadCurveData(app_config, input["name"])

	map_data := map[string]interface{}{
		"title":  curve_data.Name,
		"plot_x": GetCurveDataX(curve_data),
		"plot_y": curve_data.Values,
	}

	x_pos, err := strconv.ParseFloat(input["x"], 32)
	Check(err)

	if x_pos >= 0 {
		map_data["plot_selected_x"] = x_pos
		map_data["plot_selected_y"] = GetCurveValue(curve_data, float32(x_pos))
	}

	json_output, _ := json.Marshal(map_data)
	json_string := string(json_output)

	//log.Println("Get API Plot Result: ", json_string)

	return json_string
}

func main() {
	app_config := LoadConfig()

	actionData, err := os.ReadFile((app_config.ActionPath))
	Check(err)

	var bot Bot
	json.Unmarshal([]byte(actionData), &bot)

	engine := CreateHandlebarsEngine(app_config)

	app := CreateWebApp(engine)

	page_data_map := fiber.Map{
		"info":     "Testing 123!",
		"bot":      bot,
		"title":    "Sireus",
		"test_one": true,
		"test_two": false,
	}

	app.Post("/api/plot", func(c *fiber.Ctx) error {
		return c.SendString(GetAPIPlotData(app_config, c))
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

	app.Listen(":3000")
}
