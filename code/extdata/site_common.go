package extdata

import (
	"github.com/Knetic/govaluate"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"math"
	"strconv"
	"time"
)

func UpdateSiteBotGroups(site *appdata.Site) {
	for index, _ := range site.BotGroups {
		// Create Bots in the BotGroup from the Prometheus ExtractorKey query
		UpdateBotGroupFromPrometheus(site, index)

		// Clear all the bot variables, so our map starts fresh every time
		ClearAllBotVariables(site, index)

		// Update Bot Variables from our Queries
		UpdateBotsFromQueries(site, index)

		// Update Bot Variables from other Query Variables.  Creates Synthetic Variables.
		//NOTE(ghowland): These can be exported to Prometheus to be used in other apps, as well as Bot.ActionData
		UpdateBotsWithSyntheticVariables(site, index)

		// Update all the ActionConsiderations for each bot, so we have all the BotActionData.FinalScore values
		UpdateBotActionConsiderations(site, index)
	}
}

func UpdateBotActionConsiderations(site *appdata.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	for botIndex, bot := range botGroup.Bots {
		evalMap := GetBotEvalMapAllVariables(bot)

		for _, action := range botGroup.Actions {
			// If we don't have this ActionData yet, add it.  This will stay with the Bot for its lifetime, tracking ActiveStateTime and LastExecutionTime.
			if _, ok := site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name]; !ok {
				site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name] = appdata.BotActionData{
					ConsiderationScores: map[string]float64{},
				}
			}

			for _, consider := range action.Considerations {
				// Compile Express to be used by every bot, with their own data
				expression, err := govaluate.NewEvaluableExpression(consider.Evaluate)
				util.Check(err)

				resultInt, err := expression.Evaluate(evalMap)
				util.Check(err)

				result, err := util.ConvertInterfaceToFloat(resultInt)
				if util.Check(err) {
					// Invalidate this variable, result was invalid
					//log.Printf("Set Consideration Invalid: %s", consider.Name)
					site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationScores[consider.Name] = math.SmallestNonzeroFloat64
					continue
				}

				//log.Printf("Set Consideration Result: %s = %v", consider.Name, result)

				// Set the value.  Only valid values will exist.
				site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationScores[consider.Name] = result
			}
		}
	}
}

func ClearAllBotVariables(site *appdata.Site, botGroupIndex int) {
	for botIndex, _ := range site.BotGroups[botGroupIndex].Bots {
		site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues = map[string]float64{}
	}
}

func UpdateBotsWithSyntheticVariables(site *appdata.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	// Clear all teh Bot VariableValues

	// Create a list of names
	queryVariableNames := []string{}
	for _, variable := range botGroup.Variables {
		// Skip non-Synthetic variables
		if len(variable.Evaluate) > 0 {
			continue
		}

		queryVariableNames = append(queryVariableNames, variable.Name)
	}

	for _, variable := range botGroup.Variables {
		// Skip non-Synthetic variables
		if len(variable.Evaluate) == 0 {
			continue
		}

		// Compile Express to be used by every bot, with their own data
		expression, err := govaluate.NewEvaluableExpression(variable.Evaluate)
		util.Check(err)

		for botIndex, bot := range botGroup.Bots {
			evalMap := GetBotEvalMapOnlyQueries(bot, queryVariableNames)

			//log.Printf("Eval Map: %v", evalMap)

			resultInt, err := expression.Evaluate(evalMap)
			util.Check(err)

			result, err := util.ConvertInterfaceToFloat(resultInt)
			if util.Check(err) {
				continue // Skip this variable, it was invalid
			}

			//log.Printf("Set Synethtic Variable: %s = %v", variable.Name, result)

			// Set the value.  Only valid values will exist.
			//NOTE(ghowland): A separate test will occur to see if this bot is missing variables and cant be processed
			site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues[variable.Name] = result
		}
	}
}

func GetBotEvalMapOnlyQueries(bot appdata.Bot, queryVariableNames []string) map[string]interface{} {
	evalMap := make(map[string]interface{})

	// Build a map from this bot's variables
	for variableName, value := range bot.VariableValues {
		// Only add variables that are Query Variables, because they are known before synethetic evaluation
		if util.StringInSlice(variableName, queryVariableNames) {
			evalMap[variableName] = value
		}
	}

	return evalMap
}

func GetBotEvalMapAllVariables(bot appdata.Bot) map[string]interface{} {
	evalMap := make(map[string]interface{})

	// Build a map from this bot's variables
	for variableName, value := range bot.VariableValues {
		evalMap[variableName] = value
	}

	return evalMap
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
				// Skip variables that dont match this query, OR we have an Evaluate value, so this is an Synthetic Variable (not from Query)
				if variable.QueryName != query.Name || len(variable.Evaluate) > 0 {
					continue
				}

				//log.Printf("Bot Group: %s  Variable: %s  Key: %s == %v", botGroup.Name, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey])

				// If we have a match for this variable, next look for what Bot it matches, or it has no QueryKey we always accept it
				if len(variable.QueryKey) == 0 || (len(variable.QueryKey) > 0 && variable.QueryKeyValue == promResult.Metric[variable.QueryKey]) {
					//if variable.QueryKey == "volume" {
					//	log.Printf("Bot Group: %s   Var Bot Key: '%s'  Variable: %s  Key: %s == %v -> %v", botGroup.Name, variable.BotKey, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey], variable.QueryKeyValue == promResult.Metric[variable.QueryKey])
					//}

					for botIndex, bot := range botGroup.Bots {
						// If this Metric BotKey matches the Bot name OR the BotKey is empty, it is always accepted
						//NOTE(ghowland): Empty BotKey is used to pull data that is not specific to this Bot, but can be used as a general signal
						if promResult.Metric[variable.BotKey] == bot.Name || len(variable.BotKey) == 0 {

							//if variable.QueryName == "CPU Usage" {
							//	log.Printf("Bot Group: %s  Bot: %s   Var Bot Key: '%s'  Variable: %s  Key: %s == %v -> %v", botGroup.Name, bot.Name, variable.BotKey, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey], variable.QueryKeyValue == promResult.Metric[variable.QueryKey])
							//}

							value := math.SmallestNonzeroFloat64
							if len(promResult.Values) > 0 && len(promResult.Values[0]) > 0 {
								value, err = strconv.ParseFloat(promResult.Values[0][1].(string), 32)
								util.Check(err)
							}

							nameFormatted := util.HandlebarFormatText(variable.Name, promResult.Metric)

							site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues[nameFormatted] = float64(value)

							// If we were matching on a BotKey (normal), stop looking.  If no BotKey, do them all.
							if len(variable.BotKey) > 0 {
								break
							}
						}
					}
				}
			}
		}
	}
}
