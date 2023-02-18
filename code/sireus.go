package main

import (
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/demo"
	"github.com/ghowland/sireus/code/exporter"
	"github.com/ghowland/sireus/code/server"
	"github.com/ghowland/sireus/code/webapp"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting Sireus server...")

	// Configure the Global Server Singleton, where all the information will go
	server.Configure()

	// Run the Prometheus Exporter listener in the background
	go exporter.RunExporterListener()

	// Run the server in the background until we end the server
	go server.RunForever()

	engine := webapp.CreateHandlebarsEngine(data.SireusData.AppConfig)
	web := webapp.CreateWebApp(engine)

	// If we want to run the demo, run in the background
	if data.SireusData.AppConfig.EnableDemo {
		go demo.RunDemoForever(web)

		// Register Routes for the Demo
		webapp.RegisterRoutesDemo(web)
	}

	// Register API routes
	webapp.RegisterRoutesAPI(web)

	// Register Web routes
	webapp.RegisterRoutesWeb(web)

	// Static Files: JS, Images
	web.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./static_web"),
	}))

	_ = web.Listen(fmt.Sprintf(":%d", data.SireusData.AppConfig.WebHttpPort))
}
