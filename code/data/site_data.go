package data

type (
	// Top Level of the data structure.  Site silos all BotGroups and QueryServers, so that we can have multiple Sites
	// which are using different data sets, and should not share any data with each other.
	Site struct {
		Name          string        `json:"name"`            // Site name.  Full silo for QueryServers and BotGroups
		Info          string        `json:"info"`            // Description
		BotGroupPaths []string      `json:"bot_group_paths"` // Paths to bot_group_name.json configs
		QueryServers  []QueryServer `json:"query_servers"`   // List of QueryServers for making BotQuery requests
		BotGroups     []BotGroup    // These configure and contain ephemeral Bots which perform the Action scoring in the active States
		FreezeActions bool          // If true, no actions will be taken for this Site.  Allows control of all BotGroups.
	}
)
