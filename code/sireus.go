package main

import (
	"encoding/json"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GetActionsHtml(bot appdata.Bot) string {
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

var app_config_path = "config/config.json"

func CreateHandlebarsEngine(app_config appdata.AppConfig) *handlebars.Engine {
	// Handlebars Engine for Fiber
	engine := handlebars.New(app_config.WebPath, ".hbs")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	//// Debug will print each template that is parsed, good for debugging
	//engine.Debug(true) // Optional. Default: false

	//// Layout defines the variable name that is used to yield templates within layouts
	//engine.Layout("embed") // Optional. Default: "embed"

	raymond.RegisterHelper("botinfo", func(bot appdata.Bot) string {
		return fmt.Sprintf("%s  Actions: %d", bot.Name, len(bot.Actions))
	})

	raymond.RegisterHelper("ifconsiderlength", func(considerations []appdata.ActionConsideration, count int, options *raymond.Options) raymond.SafeString {
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

func GetCurveDataX(curve_data appdata.CurveData) []float32 {
	var x_array []float32

	for i := 0; i < len(curve_data.Values); i++ {
		x_array = append(x_array, float32(i)*0.01)
	}

	return x_array
}

func GetCurveValue(curve_data appdata.CurveData, x float32) float32 {

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
	util.Check(err)

	obj := map[string]string{}
	for k, v := range values {
		if len(v) > 0 {
			obj[k] = v[0]
		}
	}

	return obj
}

func GetAPIPlotData(app_config appdata.AppConfig, c *fiber.Ctx) string {
	input := ParseContextBody(c)
	//log.Println("Get API Plot Data: ", input)

	if input["name"] == "" {
		failure_result := map[string]interface{}{
			"_failure": "Name not found, aborting",
		}
		failure_json, _ := json.Marshal(failure_result)
		return string(failure_json)
	}

	curve_data := appdata.LoadCurveData(app_config, input["name"])

	map_data := map[string]interface{}{
		"title":  curve_data.Name,
		"plot_x": GetCurveDataX(curve_data),
		"plot_y": curve_data.Values,
	}

	x_pos, err := strconv.ParseFloat(input["x"], 32)
	util.Check(err)

	if x_pos >= 0 {
		map_data["plot_selected_x"] = x_pos
		map_data["plot_selected_y"] = GetCurveValue(curve_data, float32(x_pos))
	}

	json_output, _ := json.Marshal(map_data)
	json_string := string(json_output)

	//log.Println("Get API Plot Result: ", json_string)

	return json_string
}

func QueryPrometheus(host string, port int, query string, time_start time.Time, duration int) map[string]interface{} {
	start := time_start.UTC().Format(time.RFC3339)

	end := time_start.UTC().Add(time.Second * time.Duration(duration)).Format(time.RFC3339)

	//start := time_start.Format()

	url := fmt.Sprintf("http://%s:%d/api/v1/%s&start=%s&end=%s&step=15s", host, port, query, start, end)

	log.Print("Prom URL: ", url)

	resp, err := http.Get(url)
	util.Check(err)

	body, err := ioutil.ReadAll(resp.Body)
	util.Check(err)

	//log.Print("Prom Result: ", string(body))

	var json_result_int interface{}
	err = json.Unmarshal(body, &json_result_int)
	util.Check(err)
	json_result := json_result_int.(map[string]interface{})

	return json_result
}

func ExtractBotsFromPromData(data map[string]interface{}, bot_key string) map[string]appdata.Bot {
	//log.Print("Extra From: ", data)

	bots := make(map[string]appdata.Bot)

	result_items := data["data"].(map[string]interface{})["result"].([]interface{})

	for _, result_item := range result_items {
		item := result_item.(map[string]interface{})
		metric := item["metric"].(map[string]interface{})
		//log.Print("Item: ", metric)

		name := metric[bot_key].(string)

		_, exists := bots[name]
		if !exists {
			bots[name] = appdata.Bot{
				Name: name,
			}
		}

		//log.Print("Bot: ", name)
	}

	log.Print("Bots: ", bots)

	return bots
}

func main() {

	start_time := time.Now().Add(time.Duration(-60))
	prom_data := QueryPrometheus("localhost", 9090, "query_range?query=windows_service_status", start_time, 60)
	ExtractBotsFromPromData(prom_data, "name")

	app_config := appdata.LoadConfig(app_config_path)

	actionData, err := os.ReadFile((app_config.ActionPath))
	util.Check(err)

	var bot appdata.Bot
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
