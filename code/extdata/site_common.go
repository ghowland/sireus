package extdata

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"log"
	"time"
)

func UpdateSiteBotGroups(site *appdata.Site) {
	for index, _ := range site.BotGroups {
		UpdateBotGroupFromPrometheus(site, index) //TODO(ghowland): Test each Bot Group query extractor for source
		log.Printf("RETURN: Bots after Prom Update: %s  Count: %d", site.BotGroups[index].Name, len(site.BotGroups[index].Bots))
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

	//log.Printf("Bots after Prom Update: %s  Count: %d", botGroup.Name, len(botGroup.Bots))
}
