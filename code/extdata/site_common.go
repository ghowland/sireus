package extdata

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"math"
	"strconv"
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
	botGroup := site.BotGroups[botGroupIndex]

	// Loop over all Bot Group Queries
	for _, query := range botGroup.Queries {
		queryServer, err := appdata.GetQueryServer(*site, query.QueryServer)
		util.Check(err)

		startTime := time.Now().Add(time.Duration(-60))
		promData := QueryPrometheus(queryServer.Host, queryServer.Port, query.QueryType, query.Query, startTime, 60)

		// Loop over the Prom Results, matching Variables to Bots to save their VariableValues
		for _, promResult := range promData.Data.Result {
			// Loop through all the Variables, for every Bot.  In a Bot Group, all Bots are expected to have the same vars
			for _, variable := range botGroup.Variables {
				//log.Printf("Bot Group: %s  Variable: %s  Key: %s == %v", botGroup.Name, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey])

				// If we have a match for this variable, next look for what Bot it matches
				if variable.QueryKeyValue == promResult.Metric[variable.QueryKey] {
					//log.Printf("Match: Bot Group: %s  Variable: %s  Key: %s == %v", botGroup.Name, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey])
					for botIndex, bot := range botGroup.Bots {
						if promResult.Metric[variable.BotKey] == bot.Name {
							value := math.SmallestNonzeroFloat32
							if len(promResult.Values) > 0 && len(promResult.Values[0]) > 0 {
								value, err = strconv.ParseFloat(promResult.Values[0][1].(string), 32)
								util.Check(err)
							}

							newValue := appdata.BotVariableValue{
								Name:  variable.Name,
								Value: float32(value),
								Time:  time.Now(),
							}
							//log.Printf("Final: Bot Group: %s  Bot: %s  New Prom Value: %v", botGroup.Name, bot.Name, newValue)
							site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues = append(site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues, newValue)
							break
						}
					}
				}
			}
		}

	}
}
