package webapp

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
)

func CreateHandlebarsEngine(appConfig appdata.AppConfig) *handlebars.Engine {
	// Handlebars Engine for Fiber
	engine := handlebars.New(appConfig.WebPath, ".hbs")

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
