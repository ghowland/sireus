package main

import (
	"encoding/json"
	"fmt"
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

var action_path string = "config/action.json"
var curve_path_format string = "config/curves/%s.json"

func main() {
	actionData, err := os.ReadFile((action_path))
	Check(err)

	var bot Bot

	json.Unmarshal([]byte(actionData), &bot)

	curve_path := fmt.Sprintf(curve_path_format, bot.Actions[0].Considerations[0].CurveName)
	curveData, err := os.ReadFile((curve_path))
	Check(err)

	var curve_data CurveData
	json.Unmarshal([]byte(curveData), &curve_data)

	fmt.Println(curve_data)

	// Handlebars Engine for Fiber
	engine := handlebars.New("./web", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		//return c.SendString(GetActionsHtml(bot))
		return c.Render("index", fiber.Map{
			"info": "Testing 123!",
			"bot":  bot,
		})
	})

	app.Listen(":3000")
}
