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
func UpdateSiteBotGroups(session *data.InteractiveSession) {
	for index := range session.BotGroups {
		// Create Bots in the BotGroup from the Prometheus ExtractorKey query
		UpdateBotGroupFromPrometheus(session, &data.SireusData.Site, index)

		// Update Bot Variables from our Queries
		UpdateBotsFromQueries(session, &data.SireusData.Site, index)

		// Update Bot Variables from other Query Variables.  Creates Synthetic Variables.
		//NOTE(ghowland): These can be exported to Prometheus to be used in other apps, as well as Bot.ActionData
		UpdateBotsWithSyntheticVariables(session, index)

		// Update all the ActionConsiderations for each bot, so we have all the BotActionData.FinalScore values
		UpdateBotActionConsiderations(session, index)

		// Sort alpha, so they print consistently
		SortAllVariablesAndActions(session, index)

		// Format vars are human-readable, and we show the raw data in popups so the evaluations are clear
		CreateFormattedVariables(session, index)

		//// Execute Actions
		//ExecuteActions(session, index)
	}
}

// Execute the highest scoring action for any Bot in this Bot Group, if it is Available and meets all conditions
func ExecuteActions(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		// Lock the bot, as we are accessing the Action map
		util.LockAcquire(bot.LockKey)

		//// Look to see if we have a top-scoring action, that meets all the requirements
		//for _, actionDataPair := range bot.SortedActionData {
		//	actionData := actionDataPair.Value
		//	actionData.
		//}

		// Unlock this bot
		util.LockRelease(bot.LockKey)
	}
}

// Create formatted variables for all our Bots.  This adds human-readable strings to all the sorted Pair Lists
func CreateFormattedVariables(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		for varIndex, value := range session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues {
			variable, err := app.GetVariable(botGroup, value.Key)
			if util.CheckNoLog(err) {
				// Mark this bot as Invalid, because it is missing information
				session.BotGroups[botGroupIndex].Bots[botIndex].IsInvalid = true
				session.BotGroups[botGroupIndex].Bots[botIndex].InfoInvalid += fmt.Sprintf("Missing Variable: %s.  ", value.Key)
			}

			result := app.FormatBotVariable(variable.Format, value.Value)

			newPair := session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues[varIndex]
			newPair.Formatted = result

			session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues[varIndex] = newPair
		}
	}
}

// Sort all the Variables by name and Actions by Final Score
func SortAllVariablesAndActions(session *data.InteractiveSession, botGroupIndex int) {
	for botIndex := range session.BotGroups[botGroupIndex].Bots {
		// Cant use defer, because we are processing many in 1 action
		util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		// Sort VariableValues
		sortedVars := fixgo.SortMapStringFloat64ByKey(bot.VariableValues)
		session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues = sortedVars

		// Sort ActionData
		sortedActions := app.SortMapStringActionDataByFinalScore(bot.ActionData, false)
		session.BotGroups[botGroupIndex].Bots[botIndex].SortedActionData = sortedActions

		// Cant use defer, because we are processing many in 1 action
		util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

		//log.Printf("Bot Vars: %s  Vars: %v", bot.Name, util.PrintJson(session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues))
		//log.Printf("Bot Action Data: %s  Vars: %v", bot.Name, util.PrintJson(session.BotGroups[botGroupIndex].Bots[botIndex].SortedActionData))
	}
}

// For this BotGroup, update all the BotActionData with new ActionConsideration scores
func UpdateBotActionConsiderations(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		// Cant use defer, because we are processing many in 1 action
		util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		evalMap := GetBotEvalMapAllVariables(bot)

		for _, action := range botGroup.Actions {
			// If we don't have this ActionData yet, add it.  This will stay with the Bot for its lifetime, tracking ActiveStateTime and LastExecutionTime.
			if _, ok := session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name]; !ok {
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name] = data.BotActionData{
					ConsiderationFinalScores:  map[string]float64{},
					ConsiderationCurvedScores: map[string]float64{},
					ConsiderationRangedScores: map[string]float64{},
					ConsiderationRawScores:    map[string]float64{},
				}
			}

			for _, consider := range action.Considerations {
				// Compile Express to be used by every bot, with their own data
				expression, err := govaluate.NewEvaluableExpression(consider.Evaluate)
				util.Check(err)

				// Start assuming the data is invalid, and then mark it valid later
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationFinalScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationCurvedScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationRangedScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationRawScores[consider.Name] = math.SmallestNonzeroFloat64

				resultInt, err := expression.Evaluate(evalMap)
				if util.CheckNoLog(err) {
					// Invalidate this consideration, evaluation failed
					//log.Printf("ERROR: Evaluate failed on Eval Map data: %s   Map: %s", consider.Evaluate, util.PrintJson(evalMap))
					continue
				}

				resultRaw, err := util.ConvertInterfaceToFloat(resultInt)
				if util.CheckNoLog(err) { //TODO(ghowland): Need to handle these invalid values, so that this Bot is marked as Invalid, because the scoring cannot be done properly for every Action
					// Invalidate this consideration, result was invalid
					//log.Printf("Set Consideration Invalid: %s", consider.Name)
					continue
				}

				// Apply the Range and Curve to the Raw score
				resultRanged := util.RangeMapper(resultRaw, consider.RangeStart, consider.RangeEnd)
				curve, err := app.GetCurve(consider.CurveName)
				if util.CheckNoLog(err) {
					// Invalidate this consideration, result was invalid
					log.Printf("Set Consideration Invalid, no curve: %s   Curve missing: %s", consider.Name, consider.CurveName)
					continue
				}
				resultCurved := app.GetCurveValue(curve, resultRanged)

				//log.Printf("Set Consideration Result: %s = %v", consider.Name, resultCurved)

				considerationScore := resultCurved * consider.Weight

				// Set the value.  Only valid values will exist.
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationFinalScores[consider.Name] = considerationScore
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationCurvedScores[consider.Name] = resultCurved
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationRangedScores[consider.Name] = resultRanged
				session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name].ConsiderationRawScores[consider.Name] = resultRaw
			}

			// Get a Final Score for this Action
			calculatedScore, details := app.CalculateScore(action, session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name])
			finalScore := calculatedScore * action.Weight

			details = append(details, fmt.Sprintf("All Consider Scores: %0.2f * Action Weight: %0.2f = Final Score: %0.2f", calculatedScore, action.Weight, finalScore))

			// Copy out the ActionData struct, updated it, and assign it back into the map.
			actionData := session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name]
			actionData.FinalScore = finalScore

			allActionStatesAreActive := app.AreAllActionStatesActive(action, bot)

			// Action.WeightThreshold determines if an Action is available for possible execution
			if finalScore >= action.WeightThreshold && allActionStatesAreActive {
				if !actionData.IsAvailable {
					actionData.IsAvailable = true
					actionData.AvailableStartTime = util.GetTimeNow()
				}
			} else {
				if !allActionStatesAreActive {
					actionData.FinalScore = 0
					details = append(details, fmt.Sprintf("Setting Final Score to 0.  Missing required states: %s", util.PrintStringArrayCSV(action.RequiredStates)))
				}

				if finalScore < action.WeightThreshold {
					details = append(details, fmt.Sprintf("Final Score (%.2f) less than Action Weight Threshold (%.2f)", finalScore, action.WeightThreshold))
				}

				actionData.IsAvailable = false
				actionData.AvailableStartTime = time.UnixMilli(0)
			}

			// Details explain what happen in text, so users can better understand their results
			actionData.Details = details
			session.BotGroups[botGroupIndex].Bots[botIndex].ActionData[action.Name] = actionData
		}

		// Cant use defer, because we are processing many in 1 action
		util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
	}
}

// Update bot with Synthetic Variables.  Happens after all the Query Variables are set.  Synthetic vars can't work on each other
func UpdateBotsWithSyntheticVariables(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

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

		for botIndex := range botGroup.Bots {
			// Lock the bot
			util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

			evalMap := GetBotEvalMapOnlyQueries(session.BotGroups[botGroupIndex].Bots[botIndex], queryVariableNames)

			//log.Printf("Eval Map: %v", evalMap)

			resultInt, err := expression.Evaluate(evalMap)
			util.Check(err)

			result, err := util.ConvertInterfaceToFloat(resultInt)
			if util.Check(err) {
				util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
				continue // Skip this variable, it was invalid
			}

			//log.Printf("Set Synthetic Variable: %s = %v", variable.Name, result)

			// Set the value.  Only valid values will exist.
			//NOTE(ghowland): A separate test will occur to see if this bot is missing variables and cant be processed
			session.BotGroups[botGroupIndex].Bots[botIndex].VariableValues[variable.Name] = result

			// Unlock the bot
			util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
		}
	}
}

// Returns the map for doing the Evaluate against a Query to create our Scores.  Uses Govaluate.Evaluate()
// NOTE(ghowland): bot.AccessLock should already be locked before we come here, because we are accessing a map
func GetBotEvalMapOnlyQueries(bot data.Bot, queryVariableNames []string) map[string]interface{} {
	evalMap := make(map[string]interface{})

	// Build a map from bots variables
	for variableName, value := range bot.VariableValues {
		// Only add variables that are Query Variables, because they are known before synthetic evaluation
		if util.StringInSlice(variableName, queryVariableNames) {
			evalMap[variableName] = value
		}
	}

	return evalMap
}

// Returns the map for doing the Evaluate with a Bots VariableValues.  Uses Govaluate.Evaluate()
// NOTE(ghowland): bot.AccessLock should already be locked before we come here, because we are accessing a map
func GetBotEvalMapAllVariables(bot *data.Bot) map[string]interface{} {
	evalMap := make(map[string]interface{})

	// Build a map bot variables
	for variableName, value := range bot.VariableValues {
		evalMap[variableName] = value
	}

	return evalMap
}

// Runs Queries against Prometheus for a BotGroup
func UpdateBotGroupFromPrometheus(session *data.InteractiveSession, site *data.Site, botGroupIndex int) {
	query, err := app.GetQuery(&session.BotGroups[botGroupIndex], session.BotGroups[botGroupIndex].BotExtractor.QueryName)
	util.Check(err)

	queryResult, err := GetCachedQueryResult(session, site, query)

	//log.Printf("Extractor Query: %s", util.PrintJson(queryResult))

	extractedBots := ExtractBotsFromPromData(queryResult.PrometheusResponse, &session.BotGroups[botGroupIndex])

	//log.Printf("Extracted Bots: %s", util.PrintJson(extractedBots))

	// Find all the new bots, and add them
	//NOTE(ghowland): Removing bots is done by looking at bots that haven't had data updated past BotGroup.BotTimeoutRemove
	for _, botNew := range extractedBots {
		var foundBot bool
		for _, botCur := range session.BotGroups[botGroupIndex].Bots {
			if botCur.Name == botNew.Name {
				foundBot = true
			}
		}

		if !foundBot {
			session.BotGroups[botGroupIndex].Bots = append(session.BotGroups[botGroupIndex].Bots, botNew)
		}
	}

	// Initialize all the Bot Group states in Bot
	InitializeStates(session, botGroupIndex)
}

// Initialize all the States for this BotGroup's Bots.   They should all start at the first state value, and only move forward or reset.
func InitializeStates(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		// Clear the current states, or they grow out of control
		session.BotGroups[botGroupIndex].Bots[botIndex].StateValues = []string{}

		for _, state := range botGroup.States {
			key := fmt.Sprintf("%s.%s", state.Name, state.Labels[0])
			session.BotGroups[botGroupIndex].Bots[botIndex].StateValues = append(session.BotGroups[botGroupIndex].Bots[botIndex].StateValues, key)
		}
	}
}

// Update all the Bot VariableValues from our Queries
func UpdateBotsFromQueries(session *data.InteractiveSession, site *data.Site, botGroupIndex int) {
	botGroup := session.BotGroups[botGroupIndex]

	// Loop over all Bot Group Queries
	for _, query := range botGroup.Queries {
		// Get the cached query result, even if it is expired
		queryResult, err := GetCachedQueryResult(session, site, query)
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

					for botIndex := range botGroup.Bots {
						// Lock
						util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

						// If this Metric BotKey matches the Bot name OR the BotKey is empty, it is always accepted
						//NOTE(ghowland): Empty BotKey is used to pull data that is not specific to this Bot, but can be used as a general signal
						if promResult.Metric[variable.BotKey] == session.BotGroups[botGroupIndex].Bots[botIndex].Name || len(variable.BotKey) == 0 {

							//if variable.QueryName == "CPU Usage" {
							//	log.Printf("Bot Group: %s  Bot: %s   Var Bot Key: '%s'  Variable: %s  Key: %s == %v -> %v", botGroup.Name, bot.Name, variable.BotKey, variable.Name, variable.QueryKeyValue, promResult.Metric[variable.QueryKey], variable.QueryKeyValue == promResult.Metric[variable.QueryKey])
							//}

							value := math.SmallestNonzeroFloat64
							if len(promResult.Values) > 0 && len(promResult.Values[0]) > 0 {
								value, err = strconv.ParseFloat(promResult.Values[0][1].(string), 32)
								util.Check(err)
							}

							nameFormatted := util.HandlebarFormatText(variable.Name, promResult.Metric)

							session.BotGroups[botGroupIndex].Bots[botIndex].VariableValues[nameFormatted] = value

							// If we were matching on a BotKey (normal), stop looking.  If no BotKey, do them all.
							if len(variable.BotKey) > 0 {
								// Unlock
								util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
								break
							}
						}

						// Unlock
						util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
					}
				}
			}
		}
	}
}
