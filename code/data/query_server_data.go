package data

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
	QueryResultPool struct {
	}
)

type (
	// A single Query result
	QueryResult struct {
		QueryServer        string // Server this Query came from
		QueryType          BotQueryType
		QueryName          string             // The Query
		PrometheusResponse PrometheusResponse // The Response
	}
)
