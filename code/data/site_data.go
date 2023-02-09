package data

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
