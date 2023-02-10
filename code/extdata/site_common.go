package extdata

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/fixgo"
	"github.com/ghowland/sireus/code/util"
	"log"
	"math"
	"strconv"
	"time"
)

// Update all the BotGroups in this Site
func UpdateSiteBotGroups() {
	site := data.SireusData.Site

	for index := range site.BotGroups {
		// Create Bots in the BotGroup from the Prometheus ExtractorKey query
		UpdateBotGroupFromPrometheus(&site, index)

		// Clear all the bot variables, so our map starts fresh every time
		ClearAllBotVariables(&site, index)

		// Update Bot Variables from our Queries
		UpdateBotsFromQueries(&site, index)

		// Update Bot Variables from other Query Variables.  Creates Synthetic Variables.
		//NOTE(ghowland): These can be exported to Prometheus to be used in other apps, as well as Bot.ActionData
		UpdateBotsWithSyntheticVariables(&site, index)

		// Update all the ActionConsiderations for each bot, so we have all the BotActionData.FinalScore values
		UpdateBotActionConsiderations(&site, index)

		// Sort alpha, so they print consistently
		SortAllVariablesAndActions(&site, index)

		// Format vars are human-readable, and we show the raw data in popups so the evaluations are clear
		CreateFormattedVariables(&site, index)
	}

	// Assign the site back into the server data.  This allows atomic updates
	data.SireusData.Site = site
}

// Create formatted variables for all our Bots.  This adds human readable strings to all the sorted Pair Lists
func CreateFormattedVariables(site *data.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	for botIndex, bot := range botGroup.Bots {
		for varIndex, value := range bot.SortedVariableValues {
			variable, err := app.GetVariable(botGroup, value.Key)
			if util.CheckNoLog(err) {
				// Mark this bot as Invalid, because it is missing information
				site.BotGroups[botGroupIndex].Bots[botIndex].IsInvalid = true
				site.BotGroups[botGroupIndex].Bots[botIndex].InfoInvalid += fmt.Sprintf("Missing Variable: %s.  ", value.Key)
			}

			result := app.FormatBotVariable(variable.Format, value.Value)

			newPair := site.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues[varIndex]
			newPair.Formatted = result

			site.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues[varIndex] = newPair
		}
	}
}

// Sort all the Variables by name and Actions by Final Score
func SortAllVariablesAndActions(site *data.Site, botGroupIndex int) {
	for botIndex, bot := range site.BotGroups[botGroupIndex].Bots {
		// Sort VariableValues
		sortedVars := fixgo.SortMapStringFloat64ByKey(bot.VariableValues)
		site.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues = sortedVars

		// Sort ActionData
		sortedActions := app.SortMapStringActionDataByFinalScore(bot.ActionData, false)
		site.BotGroups[botGroupIndex].Bots[botIndex].SortedActionData = sortedActions

		//log.Printf("Bot Vars: %s  Vars: %v", bot.Name, util.PrintJson(site.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues))
		//log.Printf("Bot Action Data: %s  Vars: %v", bot.Name, util.PrintJson(site.BotGroups[botGroupIndex].Bots[botIndex].SortedActionData))
	}
}

// For this BotGroup, update all the BotActionData with new ActionConsideration scores
func UpdateBotActionConsiderations(site *data.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	for botIndex, bot := range botGroup.Bots {
		evalMap := GetBotEvalMapAllVariables(bot)

		for _, action := range botGroup.Actions {
			// If we don't have this ActionData yet, add it.  This will stay with the Bot for its lifetime, tracking ActiveStateTime and LastExecutionTime.
			if _, ok := site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name]; !ok {
				site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name] = data.BotActionData{
					ConsiderationFinalScores:     map[string]float64{},
					ConsiderationEvaluatedScores: map[string]float64{},
				}
			}

			for _, consider := range action.Considerations {
				// Compile Express to be used by every bot, with their own data
				expression, err := govaluate.NewEvaluableExpression(consider.Evaluate)
				util.Check(err)

				resultInt, err := expression.Evaluate(evalMap)
				if util.CheckNoLog(err) {
					// Invalidate this consideration, evaluation failed
					//log.Printf("ERROR: Evaluate failed on Eval Map data: %s   Map: %s", consider.Evaluate, util.PrintJson(evalMap))
					site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationFinalScores[consider.Name] = math.SmallestNonzeroFloat64
					site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationEvaluatedScores[consider.Name] = math.SmallestNonzeroFloat64
					continue
				}

				result, err := util.ConvertInterfaceToFloat(resultInt)
				if util.CheckNoLog(err) { //TODO(ghowland): Need to handle these invalid values, so that this Bot is marked as Invalid, because the scoring cannot be done properly for every Action
					// Invalidate this consideration, result was invalid
					//log.Printf("Set Consideration Invalid: %s", consider.Name)
					site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationFinalScores[consider.Name] = math.SmallestNonzeroFloat64
					site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationEvaluatedScores[consider.Name] = math.SmallestNonzeroFloat64
					continue
				}

				//log.Printf("Set Consideration Result: %s = %v", consider.Name, result)

				considerationScore := result * consider.Weight

				// Set the value.  Only valid values will exist.
				site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationFinalScores[consider.Name] = considerationScore
				site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationEvaluatedScores[consider.Name] = result
			}

			// Get a Final Score for this Action
			calculatedScore, details := app.CalculateScore(action, site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name])
			finalScore := calculatedScore * action.Weight

			// Copy out the ActionData struct, updated it, and assign it back into the map.
			actionData := site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name]
			actionData.FinalScore = finalScore

			allActionStatesAreActive := app.AreAllActionStatesActive(action, bot)

			// Action.WeightThreshold determines if an Action is available for possible execution
			if finalScore >= action.WeightThreshold && allActionStatesAreActive {
				if !actionData.IsAvailable {
					actionData.IsAvailable = true
					actionData.AvailableStartTime = time.Now()
				}
			} else {
				if !allActionStatesAreActive {
					details = append(details, fmt.Sprintf("Not all states required are active, required: %s", util.PrintStringArrayCSV(action.RequiredStates)))
				}

				if finalScore < action.WeightThreshold {
					details = append(details, fmt.Sprintf("Final Score (%.2f) did not meet Action Weight Threshold (%.2f)", finalScore, action.WeightThreshold))
				}

				actionData.IsAvailable = false
				actionData.AvailableStartTime = time.UnixMilli(0)
			}

			// Details explain what happen in text, so users can better understand their results
			actionData.Details = details
			site.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name] = actionData
		}
	}
}

// Clear all the Bot.VariableValues, so we can start fresh.  If we are missing any values, that Bot IsInvalid
func ClearAllBotVariables(site *data.Site, botGroupIndex int) {
	for botIndex := range site.BotGroups[botGroupIndex].Bots {
		site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues = map[string]float64{}
	}
}

// Update bot with Synethic Variables.  Happens after all the Query Variables are set.  Sythnetics cant work on each other
func UpdateBotsWithSyntheticVariables(site *data.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	// Clear all teh Bot VariableValues

	// Create a list of names
	var queryVariableNames []string
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

// Returns the map for doing the Evaluate against a Query to create our Scores.  Uses Govaluate.Evaluate()
func GetBotEvalMapOnlyQueries(bot data.Bot, queryVariableNames []string) map[string]interface{} {
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

// Returns the map for doing the Evaluate with a Bot's VariableValues.  Uses Govaluate.Evaluate()
func GetBotEvalMapAllVariables(bot data.Bot) map[string]interface{} {
	evalMap := make(map[string]interface{})

	// Build a map from this bot's variables
	for variableName, value := range bot.VariableValues {
		evalMap[variableName] = value
	}

	return evalMap
}

// Runs Queries against Prometheus for a BotGroup
func UpdateBotGroupFromPrometheus(site *data.Site, botGroupIndex int) {
	query, err := app.GetQuery(site.BotGroups[botGroupIndex], site.BotGroups[botGroupIndex].BotExtractor.QueryName)
	util.Check(err)

	queryServer, err := app.GetQueryServer(*site, query.QueryServer)
	util.Check(err)

	startTime := time.Now().Add(time.Duration(-60))
	promData := QueryPrometheus(queryServer.Host, queryServer.Port, query.QueryType, query.Query, startTime, 60)

	site.BotGroups[botGroupIndex].Bots = ExtractBotsFromPromData(promData, "name")

	// Initialize all the Bot Group states in Bot
	InitializeStates(site, botGroupIndex)
}

// Initialize all the States for this BotGroup's Bots.   They should all start at the first state value, and only move forward or reset.
func InitializeStates(site *data.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	for botIndex, _ := range botGroup.Bots {
		for _, state := range botGroup.States {
			key := fmt.Sprintf("%s.%s", state.Name, state.Labels[0])
			site.BotGroups[botGroupIndex].Bots[botIndex].StateValues = append(site.BotGroups[botGroupIndex].Bots[botIndex].StateValues, key)
		}
	}
}

// Update all the Bot VariableValues from our Queries
func UpdateBotsFromQueries(site *data.Site, botGroupIndex int) {
	botGroup := site.BotGroups[botGroupIndex]

	// Loop over all Bot Group Queries
	for _, query := range botGroup.Queries {
		// Get the cached query result, even if it is expired
		queryResult, err := GetCachedQueryResult(site, query, false)
		if util.CheckNoLog(err) {
			continue // Couldn't get this query, skip
		}

		// Loop over the Prom Results, matching Variables to Bots to save their VariableValues
		for _, promResult := range queryResult.PrometheusResponse.Data.Result {
			// Loop through all the Variables, for every Bot.  In a Bot Group, all Bots are expected to have the same vars
			for _, variable := range botGroup.Variables {
				// Skip variables that don't match this query, OR we have an Evaluate value, so this is a Synthetic Variable (not from Query)
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

							site.BotGroups[botGroupIndex].Bots[botIndex].VariableValues[nameFormatted] = value

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

	log.Printf("Done with initial Bot queries...")
}
