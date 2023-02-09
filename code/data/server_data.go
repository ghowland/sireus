package data

import (
	"context"
)

type (
	// SireusServerData is Singleton structure for keeping global state
	SireusServerData struct {
		AppConfig     AppConfig       // App Server configuration
		Site          Site            // For now, only 1 Site.  Later this will be dynamic
		ServerContext context.Context // Context to quickly cancel all activities
		IsQuitting    bool            // When true, this server is quitting and everything will shut down.  Controls RunUntilContextCancelled()
	}
)

var (
	// SireusData is the global data for the site, accessed by the WebApp for production and interactive ops, as well as queries
	SireusData = SireusServerData{
		IsQuitting: false,
	}
)
