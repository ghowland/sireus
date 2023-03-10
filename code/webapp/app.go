package webapp

import (
	"github.com/ghowland/sireus/code/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
)

// Initial creation of the Handlebars engine, which is passed into Fiber
func CreateHandlebarsEngine(appConfig data.AppConfig) *handlebars.Engine {
	// Handlebars Engine for Fiber
	engine := handlebars.New(appConfig.WebPath, ".hbs")

	// Reload the templates on each render, good for development
	if appConfig.ReloadTemplatesAlways {
		engine.Reload(true) // Optional. Default: false
	}

	// Debug will print each template that is parsed, good for debugging
	if appConfig.LogTemplateParsing {
		engine.Debug(true) // Optional. Default: false
	}

	// Wrap all the different helpers we will add to the handlers processing
	RegisterHandlebarsHelpers()

	return engine
}

// Create the Fiber web app, from the Handlebars engine
func CreateWebApp(engine *handlebars.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return app
}
