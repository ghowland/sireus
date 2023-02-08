package data

import (
	"context"
	"github.com/ghowland/sireus/code/appdata"
)

type (
	// SireusServerData is Singleton structure for keeping global state
	SireusServerData struct {
		AppConfig     appdata.AppConfig
		Site          appdata.Site
		ServerContext context.Context
		IsQuitting    bool // When true, this server is quitting and everything will shut down.  Controls RunUntilContextCancelled()
	}
)

var (
	// SireusData is the global data for the site, accessed by the WebApp for production and interactive ops, as well as queries
	SireusData = SireusServerData{
		IsQuitting: false,
	}
)
