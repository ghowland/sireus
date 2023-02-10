package server

import (
	"context"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/extdata"
	"github.com/ghowland/sireus/code/util"
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

	// Load for the first time too...
	LoadConfig()
}

// Load and reload config.  Will set ServerLock so this is safe to do
func LoadConfig() {
	data.SireusData.ServerLock.Lock()
	defer data.SireusData.ServerLock.Unlock()

	data.SireusData.AppConfig = app.LoadConfig(appConfigPath)

	data.SireusData.Site = app.LoadSiteConfig(data.SireusData.AppConfig)
}

// Get the global Server context, so that we can cancel everything in progress
func GetServerBackgroundContext() context.Context {
	ctx := context.Background()

	// Trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	defer func() {
		signal.Stop(channel)
		cancel()
	}()
	go func() {
		select {
		case <-channel:
			data.SireusData.IsQuitting = true //TODO(ghowland): Is this the right place for this?  Test more...
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}

// Run forever, until we stop the server
func RunForever() {
	log.Printf("Server: Run Forever: Starting (%v)", data.SireusData.IsQuitting)

	// Run until we are quitting
	for !data.SireusData.IsQuitting {
		// Run all queries that need running

		// Run all the queries that have passed their interval, or haven't been set yet
		RunAllSiteQueries(&data.SireusData.Site)

		// Update everything from the queries.  This will need time to warm up, but just let it fail in the beginning
		extdata.UpdateSiteBotGroups()

		// Pause a short time (~0.8s) to not fully spin lock the CPU ever.  This doesn't need to be more rapid
		if !data.SireusData.IsQuitting {
			time.Sleep(time.Duration(data.SireusData.AppConfig.ServerLoopDelay))
		}
	}

	log.Printf("Server: Run Forever: Stopping (%v)", data.SireusData.IsQuitting)
}

// Requests all the Queries in all the BotGroups, if they are missing or past their freshness Interval.
// Requests are not cleared, so the data will stay available for the Web App, but after the BotGroup.BotTimeoutStale
// Actions are not available.
func RunAllSiteQueries(site *data.Site) {
	for _, botGroup := range site.BotGroups {
		for _, query := range botGroup.Queries {
			// If this is already locked, then skip until the lock duration passes.  This will clear it when appropriate
			if extdata.IsQueryLocked(site, botGroup, query) {
				continue
			}

			// If we don't have this query for any reason (first time, or is over the BotQuery.Interval
			_, err := extdata.GetCachedQueryResult(site, query, true)
			if util.CheckNoLog(err) {
				go BackgroundQuery(site, query, 0)
			}
		}
	}
}

// Query in the background with a goroutine
func BackgroundQuery(site *data.Site, query data.BotQuery, interactiveUUID int64) {
	queryKey := extdata.GetQueryKey(query)

	// Set the lock, and defer to clear it when done
	extdata.QueryLockSet(site, queryKey)
	defer extdata.QueryLockClear(site, queryKey)

	// Perform the query
	queryServer, err := app.GetQueryServer(*site, query.QueryServer)
	util.Check(err)

	startTime := time.Now().Add(time.Duration(-60))
	promData := extdata.QueryPrometheus(queryServer.Host, queryServer.Port, query.QueryType, query.Query, startTime, 60)

	// Create the Query Result from
	newResult := data.QueryResult{
		QueryServer:        query.QueryServer,
		QueryType:          query.QueryType,
		Query:              query.Query,
		PrometheusResponse: promData,
	}

	extdata.StoreQueryResult(interactiveUUID, site, query, startTime, newResult)
}
