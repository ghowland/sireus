package data

import "time"

type (
	// QueryServer is where we connect to get data to populate our Bots.  example: Prometheus
	// These are stored at a Site level, so that they can be shared by all BotGroups in a Site.
	//
	// Inside a QueryServer, all QueryNames must be unique for any BotGroup, so that they can potentially be shared
	// to reduce QueryServer traffic.  Keep this in mind when creating BotGroup.Queries.
	QueryServer struct {
		ServerType          QueryServerType `json:"server_type"`
		Name                string          `json:"name"`
		Info                string          `json:"info"`
		Host                string          `json:"host"`
		Port                int             `json:"port"`
		AuthUser            string          `json:"auth_user"`
		AuthSecret          string          `json:"auth_secret"`
		DefaultStep         string          `json:"default_step"`
		DefaultDataDuration Duration        `json:"default_data_duration"`
		WebUrlFormat        string          `json:"web_url_format"`
	}
)
type (
	// QueryServerType specifies QueryServer software, defining how we make and parse Query requests
	QueryServerType int64
)

const (
	Prometheus QueryServerType = iota
)

// Format the QueryServerType for human readability
func (qst QueryServerType) String() string {
	switch qst {
	case Prometheus:
		return "Prometheus"
	}
	return "Unknown"
}

type (
	// A single Query result
	QueryResult struct {
		QueryServer        string // Server this Query came from
		QueryType          BotQueryType
		QueryName          string             // The Query
		PrometheusResponse PrometheusResponse // The Response
	}
)

type (
	// QueryResultPool is the cache for all BotGroup.Queries.  It contains normal BotQuery results from intervals, and special InteractiveUUID versions of the results, so that users can request the same query from a different time to test their Action scoring
	QueryResultPool struct {
		PoolItems []QueryResultPoolItem // These are all the items in our pool.  When we get a data request (web or internal), we get the result from here, if it exists.  New queries are run in the background and then their lateste results go here
	}
)
type (
	QueryResultPoolItem struct {
		QueryServer     string      // Server to make the query, from Site.QueryServers
		BotGroupName    string      // BotGroup where the query came from, so they namescape the QueryNames
		QueryName       string      // Name of the query in the BotGroup, query identifier
		InteractiveUUID int64       // This is 0 for normal server operation, but when a user wants to look at alternative time queries, this is set to their InteractiveUUID
		TimeRequested   time.Time   // Time the BotQuery was requested
		TimeReceived    time.Time   // Time the Response was received
		Result          QueryResult // Response from the QueryServer
		IsValid         bool        // Is the response valid?  If false, it can't be used
	}
)
