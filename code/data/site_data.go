package data

type (
	// Top Level of the data structure.  Site silos all BotGroups and QueryServers, so that we can have multiple Sites
	// which are using different data sets, and should not share any data with each other.
	Site struct {
		Name                    string                 `json:"name"`            // Site name.  Full silo for QueryServers and BotGroups
		Info                    string                 `json:"info"`            // Description
		BotGroupPaths           []string               `json:"bot_group_paths"` // Paths to bot_group_name.json configs
		QueryServers            []QueryServer          `json:"query_servers"`   // List of QueryServers for making BotQuery requests
		FreezeActions           bool                   // If true, no actions will be taken for this Site.  Allows control of all BotGroups Action execution.
		QueryResultCache        QueryResultPool        // Per Site, we cache all the BotQuery results here.  Per normal server operation, and per InteractiveSession
		InteractiveSessionCache InteractiveSessionPool // Per Site, we track web app InteractiveSession data to allow users to make changes and see how they alter the Action scoring.  Sites silo everything, so it would be an anti-feature to allow InteractiveSession data to cross Site boundarie
		LoadedBotGroups         []BotGroup             // These are just JSON loaded values to be cloned for the InteractiveSesssion.BotGroups, which contain Bots which perform the Action scoring in the active States
		ProductionControl       InteractiveControl     // This is the config loaded production (UUID=0) version of InteractiveControl.  Storing it here means it doesn't have to keep being generated when needed.
	}
)
