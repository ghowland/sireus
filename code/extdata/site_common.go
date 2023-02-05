package extdata

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"time"
)

func UpdateSiteBotGroups(site *appdata.Site) {
	for index, _ := range site.BotGroups {
		UpdateBotGroupFromPrometheus(site, index)

		UpdateBotsFromQueries(site, index)
	}
}

func UpdateBotGroupFromPrometheus(site *appdata.Site, botGroupIndex int) {
	query, err := appdata.GetQuery(site.BotGroups[botGroupIndex], site.BotGroups[botGroupIndex].BotExtractor.QueryName)
	util.Check(err)

	queryServer, err := appdata.GetQueryServer(*site, query.QueryServer)
	util.Check(err)

	startTime := time.Now().Add(time.Duration(-60))
	promData := QueryPrometheus(queryServer.Host, queryServer.Port, query.QueryType, query.Query, startTime, 60)

	site.BotGroups[botGroupIndex].Bots = ExtractBotsFromPromData(promData, "name")
}

func UpdateBotsFromQueries(site *appdata.Site, botGroupIndex int) {
	/*
		botGroup := site.BotGroups[botGroupIndex]

		// Loop over all Bot Group Queries
		for _, query := range botGroup.Queries {
			queryServer, err := appdata.GetQueryServer(*site, query.QueryServer)
			util.Check(err)

			startTime := time.Now().Add(time.Duration(-60))
			promData := QueryPrometheus(queryServer.Host, queryServer.Port, query.QueryType, query.Query, startTime, 60)

			//for _,

			// Loop through all the Variables, for every Bot.  In a Bot Group, all Bots are expected to have the same vars
			for _, variable := range botGroup.Variables {
				for _, bot := range botGroup.Bots {

				}
			}

		}
	*/
}
