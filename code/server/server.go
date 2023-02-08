package server

import (
	"context"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/extdata"
	"log"
	"os"
	"os/signal"
	"time"
)

// TODO(ghowland): Handle config file as CLI arg
var appConfigPath = "config/config.json"

func Configure() {
	data.SireusData.IsQuitting = false

	data.SireusData.ServerContext = GetServerBackgroundContext()

	data.SireusData.AppConfig = appdata.LoadConfig(appConfigPath)

	data.SireusData.Site = appdata.LoadSiteConfig(data.SireusData.AppConfig)
}

// Get the global Server context, so that we can cancel everything in progress
func GetServerBackgroundContext() context.Context {
	ctx := context.Background()

	// Trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	defer func() {
		data.SireusData.IsQuitting = true
		signal.Stop(channel)
		cancel()
	}()
	go func() {
		select {
		case <-channel:
			data.SireusData.IsQuitting = true
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}

// Run forever, until we stop the server
func RunForever() {
	log.Printf("Server: Run Forever: Starting (%v)", data.SireusData.IsQuitting)

	//for !data.SireusData.IsQuitting {
	for true {
		extdata.UpdateSiteBotGroups()

		// Simplest method of delay initially, just keep updating everything.
		// Next will move this to following the individual query configs
		time.Sleep(5 * time.Second)
	}

	log.Printf("Server: Run Forever: Stopping (%v)", data.SireusData.IsQuitting)
}
