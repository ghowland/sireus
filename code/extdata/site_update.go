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
	"sort"
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
		//NOTE(ghowland): These can be exported to Prometheus to be used in other apps, as well as Bot.ConditionData
		UpdateBotsWithSyntheticVariables(session, index)

		// Export Metrics on Variables marked for export
		ExportMetricsOnVariables(session, index)

		// Update all the ConditionConsiderations for each bot, so we have all the BotConditionData.FinalScore values
		UpdateBotConditionConsiderations(session, index)

		// Sort alpha, so they print consistently
		SortAllVariablesAndConditions(session, index)

		// Execute Conditions (lock and delay testing inside)
		executedConditions := ExecuteBotGroupConditions(session, index)

		// If we executed conditions, we need to make sure things are updated and sorted again, because they have changed
		if executedConditions {
			// Repeat this, to ensure things that are now Inactive after a state change from Executing Conditions
			UpdateBotConditionConsiderations(session, index)
			SortAllVariablesAndConditions(session, index)
		}

		// Format vars are human-readable, and we show the raw data in popups so the evaluations are clear
		CreateFormattedVariables(session, index)
	}
}

// Execute the highest scoring condition for any Bot in this Bot Group, if it is Available and meets all conditions
func ExecuteBotGroupConditions(session *data.InteractiveSession, botGroupIndex int) bool {
	botGroup := &session.BotGroups[botGroupIndex]

	executedConditions := false

	for botIndex := range botGroup.Bots {
		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		// Lock the bot, as we are accessing the Condition map
		util.LockAcquire(bot.LockKey)

		// Take the top scoring item only and see if it is available and meets any additional requirements
		if bot.SortedConditionData.Len() > 0 {
			conditionDataName := bot.SortedConditionData[0].Key
			conditionData := bot.SortedConditionData[0].Value

			condition, err := app.GetCondition(botGroup, conditionDataName)
			if util.Check(err) {
				log.Printf("Missing Condition: %s   Bot Group: %s  Bot: %s", conditionDataName, botGroup.Name, bot.Name)
				continue
			}

			// If the condition is available, and the final score is over the threshold, test next steps
			if conditionData.IsAvailable && conditionData.FinalScore > botGroup.ConditionThreshold {
				timeAvailable := time.Now().Sub(conditionData.AvailableStartTime)

				// If we have been available for long enough, this should be the final check, we can execute this Condition
				if timeAvailable.Seconds() > time.Duration(condition.RequiredAvailable).Seconds() {
					ExecuteBotCondition(session, botGroup, bot, condition, conditionData)
					executedConditions = true
				}
			}
		}

		// Unlock this bot
		util.LockRelease(bot.LockKey)
	}

	return executedConditions
}

func ExecuteBotCondition(session *data.InteractiveSession, botGroup *data.BotGroup, bot *data.Bot, condition data.Condition, conditionData data.BotConditionData) {
	// Lock this session for execution.  We want to be able to delay them, and ensure they aren't racing, because HTTP requests trigger this
	util.LockAcquire(fmt.Sprintf("session.%d.execute_condition.%s", session.UUID, condition.Name))
	defer util.LockRelease(fmt.Sprintf("session.%d.execute_condition.%s", session.UUID, condition.Name))

	// Get the last time this Condition executed, so we can enforce a repeat execution delay
	conditionLastExecuteTime, err := app.GetConditionLastExecuteTime(session, botGroup, bot, condition, condition.ExecuteRepeatDelay)

	// Return early if we executed within the delay threshold.  In this case, err means it wasn't executed, so we will perform the execution.  err is not a failure case here
	if !util.Check(err) && util.GetTimeNow().Sub(conditionLastExecuteTime) < time.Duration(condition.ExecuteRepeatDelay) {
		//log.Printf(fmt.Sprintf("Session Bot Execute Conditions returning early because called too soon: %d  Last: %v  Cur: %v", session.UUID, util.GetTimeNow(), util.GetTimeNow()))
		return
	}

	// Create the Condition Command Result which will go into the Bot Command History
	commandResult := data.ConditionCommandResult{
		BotGroupName:  botGroup.Name,
		BotName:       bot.Name,
		ConditionName: condition.Name,
		Started:       util.GetTimeNow(),
		Score:         conditionData.FinalScore,
		StatesBefore:  util.CopyStringSlice(bot.StateValues),
	}

	// Format the CommandLog so we have a rich version
	formatMap := map[string]interface{}{
		"botGroup":         botGroup,
		"bot":              bot,
		"condition":        condition,
		"conditionCommand": condition.Command,
		"conditionData":    conditionData,
		"appConfig":        data.SireusData.AppConfig,
	}
	commandResult.CommandLog = util.HandlebarFormatData(condition.Command.LogFormat, formatMap)

	// Set the Lock Timers
	app.SetAllConditionLockTimers(condition, botGroup, condition.Command.LockTimerDuration)

	// Update the states
	err = app.SetBotStates(botGroup, bot, condition.Command.SetBotStates)
	if util.Check(err) {
		log.Printf("Aborting condition execution, can't set state: Invalid configuration, states were not successfully updated and may be out of sync with each other now: Bot Group: %s  Bot: %s  Condition: %s  Error: %s", botGroup.Name, bot.Name, condition.Name, err.Error())
		return
	}

	// Reset any states required
	for _, resetState := range condition.Command.ResetBotStates {
		err := app.ResetBotState(botGroup, bot, resetState)
		if util.Check(err) {
			log.Printf("Aborting condition execution, can't reset state: Invalid configuration, states were not successfully updated and may be out of sync with each other now: Bot Group: %s  Bot: %s  Condition: %s  Error: %s", botGroup.Name, bot.Name, condition.Name, err.Error())
			return
		}
	}

	// Save the States after our changes
	commandResult.StatesAfter = util.CopyStringSlice(bot.StateValues)

	// Execute command
	//TODO(ghowland): Move this to the Sireus Client, so it can be run in different locations to get different access
	if condition.Command.Type == 1 || condition.Command.Type == 2 {
		url := util.HandlebarFormatData(condition.Command.Content, formatMap)
		body, err := util.HttpGet(url)
		if util.Check(err) {
			commandResult.ResultContent = fmt.Sprintf("Error: %s", err.Error())
			//log.Printf(fmt.Sprintf("%s: %s: %s == Error: %s", botGroup.Name, bot.Name, url, err.Error()))
		} else {
			commandResult.ResultContent = body
			//log.Printf(fmt.Sprintf("%s: %s: %s == %s", botGroup.Name, bot.Name, url, body))
		}
	}

	// Mark our completion time
	commandResult.Finished = util.GetTimeNow()

	// Append the Command Result to the Bots Command History
	bot.CommandHistory = append(bot.CommandHistory, commandResult)

	// Increment the Metric Counter, that we executed this Condition's Command
	app.AddToMetricCounter("sireus_execute_condition", 1, "A Condition met all the requirements and had the highest score, so was executed", app.GetMetricLabelsAndInfo_Condition(botGroup, bot, condition))
}

// Create formatted variables for all our Bots.  This adds human-readable strings to all the sorted Pair Lists
func CreateFormattedVariables(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		for varIndex, value := range session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues {
			variable, err := app.GetVariable(botGroup, value.Key)
			if util.Check(err) {
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

func ExportMetricsOnVariables(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	// For all the BotGroup variables that are marked for Export...
	for _, varData := range botGroup.Variables {
		if varData.Export {
			// Export them for every bot that has them
			for botIndex := range botGroup.Bots {
				bot := &session.BotGroups[botGroupIndex].Bots[botIndex]
				value, ok := bot.VariableValues[varData.Name]
				if !ok {
					continue
				}

				app.SetMetricGauge("sireus_variable", value, "A Bot variable marked for exporting, probably synthesized", app.GetMetricLabelsAndInfo_BotVariable(botGroup, bot, varData.Name))
			}
		}
	}

}

// Sort all the Variables by name and Conditions by Final Score
func SortAllVariablesAndConditions(session *data.InteractiveSession, botGroupIndex int) {
	for botIndex := range session.BotGroups[botGroupIndex].Bots {
		// Cant use defer, because we are processing many in 1 condition
		util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		// Sort VariableValues
		sortedVars := fixgo.SortMapStringFloat64ByKey(bot.VariableValues)
		session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues = sortedVars

		// Sort ConditionData
		sortedConditions := app.SortMapStringConditionDataByFinalScore(bot.ConditionData, false)
		session.BotGroups[botGroupIndex].Bots[botIndex].SortedConditionData = sortedConditions

		// Cant use defer, because we are processing many in 1 condition
		util.LockRelease(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

		//log.Printf("Bot Vars: %s  Vars: %v", bot.Name, util.PrintJson(session.BotGroups[botGroupIndex].Bots[botIndex].SortedVariableValues))
		//log.Printf("Bot Condition Data: %s  Vars: %v", bot.Name, util.PrintJson(session.BotGroups[botGroupIndex].Bots[botIndex].SortedConditionData))
	}
}

// For this BotGroup, update all the BotConditionData with new ConditionConsideration scores
func UpdateBotConditionConsiderations(session *data.InteractiveSession, botGroupIndex int) {
	botGroup := &session.BotGroups[botGroupIndex]

	for botIndex := range botGroup.Bots {
		// Cant use defer, because we are processing many in 1 condition
		util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)
		bot := &session.BotGroups[botGroupIndex].Bots[botIndex]

		evalMap := GetBotEvalMapAllVariables(bot)

		for _, condition := range botGroup.Conditions {
			// If we don't have this ConditionData yet, add it.  This will stay with the Bot for its lifetime, tracking ActiveStateTime and LastExecutionTime.
			if _, ok := session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name]; !ok {
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name] = data.BotConditionData{
					ConsiderationFinalScores:  map[string]float64{},
					ConsiderationCurvedScores: map[string]float64{},
					ConsiderationRangedScores: map[string]float64{},
					ConsiderationRawScores:    map[string]float64{},
				}
			}

			for _, consider := range condition.Considerations {
				// Compile Express to be used by every bot, with their own data
				expression, err := govaluate.NewEvaluableExpression(consider.Evaluate)
				util.CheckLog(err)

				// Start assuming the data is invalid, and then mark it valid later
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationFinalScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationCurvedScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationRangedScores[consider.Name] = math.SmallestNonzeroFloat64
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationRawScores[consider.Name] = math.SmallestNonzeroFloat64

				resultInt, err := expression.Evaluate(evalMap)
				if util.Check(err) {
					// Invalidate this consideration, evaluation failed
					//log.Printf("ERROR: Evaluate failed on Eval Map data: %s   Map: %s", consider.Evaluate, util.PrintJson(evalMap))
					continue
				}

				resultRaw, err := util.ConvertInterfaceToFloat(resultInt)
				if util.Check(err) { //TODO(ghowland): Need to handle these invalid values, so that this Bot is marked as Invalid, because the scoring cannot be done properly for every Condition
					// Invalidate this consideration, result was invalid
					//log.Printf("Set Consideration Invalid: %s", consider.Name)
					continue
				}

				// Apply the Range and Curve to the Raw score
				resultRanged := util.RangeMapper(resultRaw, consider.RangeStart, consider.RangeEnd)
				curve, err := app.GetCurve(consider.CurveName)
				if util.Check(err) {
					// Invalidate this consideration, result was invalid
					log.Printf("Set Consideration Invalid, no curve: %s   Curve missing: %s", consider.Name, consider.CurveName)
					continue
				}
				resultCurved := app.GetCurveValue(curve, resultRanged)

				//log.Printf("Set Consideration Result: %s = %v", consider.Name, resultCurved)

				considerationScore := resultCurved * consider.Weight

				// Set the value.  Only valid values will exist.
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationFinalScores[consider.Name] = considerationScore
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationCurvedScores[consider.Name] = resultCurved
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationRangedScores[consider.Name] = resultRanged
				session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name].ConsiderationRawScores[consider.Name] = resultRaw
			}

			// Get a Final Score for this Condition
			calculatedScore, details := app.CalculateScore(condition, session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name])
			finalScore := calculatedScore * condition.Weight

			details = append(details, fmt.Sprintf("All Consider Scores: %0.2f * Condition Weight: %0.2f = Final Score: %0.2f", calculatedScore, condition.Weight, finalScore))

			// Copy out the ConditionData struct, updated it, and assign it back into the map.
			conditionData := session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name]
			conditionData.FinalScore = finalScore

			allConditionStatesAreActive := app.AreAllConditionStatesActive(condition, bot)

			allConditionRequiredLocksTimersAvailable := app.AreAllConditionLockTimersAvailable(condition, botGroup)

			// Condition.WeightThreshold determines if a Condition is available for possible execution
			if finalScore >= condition.WeightThreshold && allConditionStatesAreActive && allConditionRequiredLocksTimersAvailable {
				if !conditionData.IsAvailable {
					conditionData.IsAvailable = true
					conditionData.AvailableStartTime = util.GetTimeNow()
				}
			} else {
				if !allConditionStatesAreActive {
					conditionData.FinalScore = 0
					details = append(details, fmt.Sprintf("Setting Final Score to 0.  Missing required states: %s", util.PrintStringArrayCSV(condition.RequiredStates)))
				}

				if !allConditionRequiredLocksTimersAvailable {
					details = append(details, fmt.Sprintf("Not available.  Missing required Lock Timers: %s", util.PrintStringArrayCSV(condition.RequiredLockTimers)))
				}

				if finalScore < condition.WeightThreshold {
					details = append(details, fmt.Sprintf("Final Score (%.2f) less tha Condition Weight Threshold (%.2f)", finalScore, condition.WeightThreshold))
				}

				conditionData.IsAvailable = false
				conditionData.AvailableStartTime = time.UnixMilli(0)
			}

			// Details explain what happen in text, so users can better understand their results
			conditionData.Details = details
			session.BotGroups[botGroupIndex].Bots[botIndex].ConditionData[condition.Name] = conditionData
		}

		// Cant use defer, because we are processing many in 1 condition
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
		util.CheckLog(err)

		for botIndex := range botGroup.Bots {
			// Lock the bot
			util.LockAcquire(session.BotGroups[botGroupIndex].Bots[botIndex].LockKey)

			evalMap := GetBotEvalMapOnlyQueries(session.BotGroups[botGroupIndex].Bots[botIndex], queryVariableNames)

			//log.Printf("Eval Map: %v", evalMap)

			resultInt, err := expression.Evaluate(evalMap)
			util.CheckLog(err)

			result, err := util.ConvertInterfaceToFloat(resultInt)
			if util.CheckLog(err) {
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
		if util.StringInSlice(queryVariableNames, variableName) {
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
	util.CheckLog(err)

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
			// Initialize all the Bot Group states in Bot
			InitializeBotStates(&session.BotGroups[botGroupIndex], &botNew)

			// Add it into the botGroup slice
			session.BotGroups[botGroupIndex].Bots = append(session.BotGroups[botGroupIndex].Bots, botNew)
		}
	}
}

// Initialize all the States for this BotGroups Bots.   They should all start at the first state value, and only move forward or reset.
func InitializeBotStates(botGroup *data.BotGroup, bot *data.Bot) {
	// Clear the current states, or they grow out of control
	bot.StateValues = []string{}

	for _, state := range botGroup.States {
		key := fmt.Sprintf("%s.%s", state.Name, state.Labels[0])
		bot.StateValues = append(bot.StateValues, key)
	}

	sort.Strings(bot.StateValues)
}

// Update all the Bot VariableValues from our Queries
func UpdateBotsFromQueries(session *data.InteractiveSession, site *data.Site, botGroupIndex int) {
	botGroup := session.BotGroups[botGroupIndex]

	// Loop over all Bot Group Queries
	for _, query := range botGroup.Queries {
		// Get the cached query result, even if it is expired
		queryResult, err := GetCachedQueryResult(session, site, query)
		if util.Check(err) {
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
								util.CheckLog(err)
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
