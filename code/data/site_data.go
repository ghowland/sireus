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
	// Top Level of the data structure.  Site silos all BotGroups and QueryServers, so that we can have multiple Sites
	// which are using different data sets, and should not share any data with each other.
	Site struct {
		Name          string        `json:"name"`
		Info          string        `json:"info"`
		BotGroupPaths []string      `json:"bot_group_paths"`
		QueryServers  []QueryServer `json:"query_servers"`
		BotGroups     []BotGroup
		FreezeActions bool // If true, no actions will be taken for this Site.  Allows control of all BotGroups.
	}
)
