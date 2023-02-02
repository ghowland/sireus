package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

var action_path string = "config/action.json"

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

func main() {
	actionData, err := os.ReadFile((action_path))
	Check(err)

	var bot Bot

	json.Unmarshal([]byte(actionData), &bot)

	fmt.Println(bot)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(GetActionsHtml(bot))
	})

	app.Listen(":3000")
}
