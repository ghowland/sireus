package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/fixgo"
	"github.com/ghowland/sireus/code/util"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Load the BotGroup config from a path
func LoadBotGroupConfig(path string) data.BotGroup {
	botGroupData, err := os.ReadFile(path)
	util.CheckPanic(err)

	var botGroup data.BotGroup
	err = json.Unmarshal(botGroupData, &botGroup)
	util.CheckPanic(err)

	return botGroup
}

// Load our Site config for a path
func LoadSiteConfig(appConfig data.AppConfig) data.Site {
	siteData, err := os.ReadFile(appConfig.SiteConfigPath)
	util.CheckPanic(err)

	var site data.Site
	err = json.Unmarshal(siteData, &site)
	util.CheckPanic(err)

	// Initialize data that isn't auto-initialized or loaded from JSON
	site.ProductionControl = GetProductionInteractiveControl()
	site.InteractiveSessionCache.Sessions = make(map[data.SessionUUID]data.InteractiveSession)
	site.QueryResultCache = data.QueryResultPool{
		PoolItems:  make(map[string]data.QueryResultPoolItem),
		QueryLocks: make(map[string]time.Time),
	}

	// Load all our Bot Groups.  We keep these cached for cloning, so we don't have to parse JSON all the time, but put nothing dynamic into them
	for _, botGroupPath := range site.BotGroupPaths {
		botGroup := LoadBotGroupConfig(botGroupPath)
		botGroup.LockKey = fmt.Sprintf("%s.%s", site.Name, botGroup.Name)

		site.LoadedBotGroups = append(site.LoadedBotGroups, botGroup)
	}

	return site
}

// Takes an InteractiveControl struct, and creates a InteractiveSession, which is used everywhere and contains live BotGroups
func GetInteractiveSession(interactiveControl data.InteractiveControl, site *data.Site) data.InteractiveSession {
	site.InteractiveSessionCache.AccessLock.Lock()
	defer site.InteractiveSessionCache.AccessLock.Unlock()

	// Get an existing session for this UUID
	session, ok := site.InteractiveSessionCache.Sessions[interactiveControl.SessionUUID]
	if !ok {
		// Couldn't find it, so create one
		session = data.InteractiveSession{
			UUID:      interactiveControl.SessionUUID,
			BotGroups: site.LoadedBotGroups,
		}
		site.InteractiveSessionCache.Sessions[interactiveControl.SessionUUID] = session
	}

	// Always update these values
	session.TimeRequested = util.GetTimeNow()
	session.QueryStartTime = time.UnixMilli(int64(interactiveControl.QueryStartTime))
	session.QueryDuration = data.Duration(interactiveControl.QueryDuration)
	session.QueryScrubTime = time.UnixMilli(int64(interactiveControl.QueryScrubTime))

	// If this is a Production session, we invalidate cache on interval and don't worry about query mismatch
	if interactiveControl.SessionUUID == 0 {
		session.IgnoreCacheQueryMismatch = false
		session.IgnoreCacheOverInterval = true
	} else {
		// Else, this is an Interactive session, so we don't ignore old queries.  We just want them to match time
		session.IgnoreCacheOverInterval = false
		session.IgnoreCacheQueryMismatch = true
	}

	// We modified the session, put it back into the session map
	site.InteractiveSessionCache.Sessions[interactiveControl.SessionUUID] = session

	return session
}

// GetProductionInteractiveControl returns a SessionUUID==0 data set for production data.
// TODO(ghowland): These should be altered by AppConfig
func GetProductionInteractiveControl() data.InteractiveControl {
	interactiveControl := data.InteractiveControl{
		SessionUUID:            0,
		UseInteractiveSession:  false,
		UseInteractiveOverride: false,
		QueryStartTime:         float64(util.GetTimeNow().UnixMilli()),
		QueryDuration:          60 * 1000000000,
		QueryScrubTime:         float64(util.GetTimeNow().UnixMilli()),
	}

	return interactiveControl
}

// Returns a QueryServer, scope is per Site
func GetQueryServer(site *data.Site, name string) (data.QueryServer, error) {
	for _, queryServer := range site.QueryServers {
		if queryServer.Name == name {
			return queryServer, nil
		}
	}

	return data.QueryServer{}, errors.New(fmt.Sprintf("Query Server missing: %s", name))
}

// Gets a query, scope per BotGroup
func GetQuery(botGroup *data.BotGroup, queryName string) (data.BotQuery, error) {
	for _, query := range botGroup.Queries {
		if query.Name == queryName {
			return query, nil
		}
	}
	return data.BotQuery{}, errors.New(fmt.Sprintf("Bot Group: %s  Query missing: %s", botGroup.Name, queryName))
}

// Get a Condition from a BotGroup, by name
func GetCondition(botGroup *data.BotGroup, actionName string) (data.Condition, error) {
	for _, action := range botGroup.Conditions {
		if action.Name == actionName {
			return action, nil
		}
	}
	return data.Condition{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Condition: %s", botGroup.Name, actionName))
}

// Get a Variable defintion from BotGroup, by name.  Not the Variable Value, which is stored in Bot.
func GetVariable(botGroup *data.BotGroup, varName string) (data.BotVariable, error) {
	for _, variable := range botGroup.Variables {
		if variable.Name == varName {
			return variable, nil
		}
	}
	return data.BotVariable{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Variable: %s", botGroup.Name, varName))
}

// Get a ConditionConsideration from a Condition, by name
func GetConditionConsideration(action data.Condition, considerName string) (data.ConditionConsideration, error) {
	for _, consider := range action.Considerations {
		if consider.Name == considerName {
			return consider, nil
		}
	}
	return data.ConditionConsideration{}, errors.New(fmt.Sprintf("Missing Consideration: %s", considerName))
}

func GetConditionLastExecuteTime(session *data.InteractiveSession, botGroup *data.BotGroup, bot *data.Bot, action data.Condition, stopLookingAfter data.Duration) (time.Time, error) {
	// If we have no history, then this hasn't been run before
	if len(bot.CommandHistory) == 0 {
		return time.Time{}, errors.New(fmt.Sprintf("This command has never been run: %d  Bot Group: %s  Bot: %s  Condition: %s", session.UUID, botGroup.Name, bot.Name, action.Name))
	}

	// Loop backwards over the command history, looking for this Condition
	for i := len(bot.CommandHistory) - 1; i >= 0; i-- {
		commandResult := bot.CommandHistory[i]

		// If this is a match, return successfully
		if commandResult.ConditionName == action.Name {
			return commandResult.Started, nil
		}

		// If this is past our time to stop looking, then this hasn't been run in the time we care about
		if util.GetTimeNow().Sub(commandResult.Started) > time.Duration(stopLookingAfter) {
			return time.Time{}, errors.New(fmt.Sprintf("This command has not been run since the timeout: %d  Bot Group: %s  Bot: %s  Condition: %s  Timeout: %v", session.UUID, botGroup.Name, bot.Name, action.Name, stopLookingAfter))
		}
	}

	return time.Time{}, errors.New(fmt.Sprintf("This command has not been in the entire command history: %d  Bot Group: %s  Bot: %s  Condition: %s  Timeout: %v", session.UUID, botGroup.Name, bot.Name, action.Name, stopLookingAfter))
}

// For a given Condition, does this Bot have all the RequiredStates active?
func AreAllConditionStatesActive(action data.Condition, bot *data.Bot) bool {
	for _, state := range action.RequiredStates {
		if !util.StringInSlice(bot.StateValues, state) {
			return false
		}
	}
	return true
}

// For a given Condition, does this Bot Group have all the Lock Timers available to be locked?
func AreAllConditionLockTimersAvailable(action data.Condition, botGroup *data.BotGroup) bool {
	for _, lockTimerName := range action.RequiredLockTimers {
		lockTimer, err := GetLockTimer(botGroup, lockTimerName)
		if util.Check(err) {
			log.Printf("Missing Lock Timer: %s  Invalid configuration, will never activate Condition: %s  Bot Group: %s", lockTimerName, action.Name, botGroup.Name)
			return false
		}

		// If this lock timer is active, return false.  Not all lock timers are available
		if lockTimer.IsActive {
			return false
		}
	}
	return true
}

// When executing a Condition, we will set all the Lock Timers that Condition required, for the duration specified in the ConditionCommand
func SetAllConditionLockTimers(action data.Condition, botGroup *data.BotGroup, duration data.Duration) {
	for _, lockTimerName := range action.RequiredLockTimers {
		SetLockTimer(botGroup, lockTimerName, duration)
	}
}

func ResetBotState(botGroup *data.BotGroup, bot *data.Bot, stateBase string) error {
	// Get our state data, so we can get the first label
	stateData, err := GetBotForwardSequenceState(botGroup, stateBase)
	if util.Check(err) {
		return err
	}

	// Get the current state index, so we can remove it
	currentStateName, _, err := GetBotCurrentStateAndIndex(botGroup, bot, stateBase)
	if util.Check(err) {
		return err
	}

	// Remove the current state
	bot.StateValues, _ = util.StringSliceRemoveString(bot.StateValues, currentStateName)

	// Add the default state
	key := fmt.Sprintf("%s.%s", stateData.Name, stateData.Labels[0])
	bot.StateValues = append(bot.StateValues, key)

	// Sort so they are in a consistent order
	sort.Strings(bot.StateValues)

	return nil
}

// When executing a Condition, we want to update the Bots States, to move it forward
func SetBotStates(botGroup *data.BotGroup, bot *data.Bot, setStates []string) error {
	// Update the states
	for _, state := range setStates {
		stateBase := ""
		stateLabel := ""
		advanceOnly := false

		// Get the State Base and Target, so we can remove any existing states that are prefixed with this
		if strings.Contains(state, ".") {
			stateSplit := strings.SplitN(state, ".", 2)
			stateBase = stateSplit[0]
			stateLabel = stateSplit[1]
		} else {
			// Else, we only have a base, so we will advance the state forward, instead of setting a target
			stateBase = stateLabel
			advanceOnly = true
		}

		// Get our current state index, for this stateBase
		currentState, currentStateIndex, err := GetBotCurrentStateAndIndex(botGroup, bot, stateBase)
		if util.Check(err) {
			log.Printf(err.Error())
			return err
		}

		// Setting the state by name is the normal case, allows skipping steps.  Also, if you want a new state inserted that's a single config change
		if !advanceOnly {
			// Get the target state index, as it has to be higher than our current index, or that's invalid, and we can't execute
			targetIndex, err := GetStateIndex(botGroup, state)
			if util.Check(err) {
				return err
			}
			if targetIndex < currentStateIndex {
				return errors.New(fmt.Sprintf("Trying to set a state to an earlier label, which is not allowed.  They can only progress forward.  Target: %s (%d)  Current: %s (%d)  ", state, targetIndex, currentState, currentStateIndex))
			}

			// Remove the current state, we already have its index
			bot.StateValues, err = util.StringSliceRemoveString(bot.StateValues, currentState)
			if util.Check(err) {
				return err
			}

			// Add the new state
			bot.StateValues = append(bot.StateValues, state)
		} else {
			stateData, err := GetBotForwardSequenceState(botGroup, stateBase)
			if util.Check(err) {
				return err
			}

			// Advance Only.  We will keep incrementing the state, until we reach the end, and then we will stay there until reset
			lastIndex := len(stateData.Labels) - 1

			if currentStateIndex < lastIndex {
				nextState := stateData.Labels[currentStateIndex+1]

				// Remove the current state, we already have its index
				bot.StateValues, err = util.StringSliceRemoveString(bot.StateValues, currentState)
				if util.Check(err) {
					return err
				}

				// Add the new state
				bot.StateValues = append(bot.StateValues, nextState)
			} else {
				// We don't need to advance anymore, but this is not an error, we are waiting until we are reset back to 0
			}
		}
	}

	// Sort the Bot.StateValues, so they are consistent when reading them
	sort.Strings(bot.StateValues)

	return nil
}

func GetBotForwardSequenceState(botGroup *data.BotGroup, name string) (data.BotForwardSequenceState, error) {
	for _, state := range botGroup.States {
		if state.Name == name {
			return state, nil
		}
	}

	return data.BotForwardSequenceState{}, errors.New(fmt.Sprintf("Missing Bot Forward Sequence State: %s  Bot Group: %s", name, botGroup.Name))
}

// Returns the index of the State currently for this Bot, with the stateBase (BotForwardSequenceState.Name)
func GetBotCurrentStateAndIndex(botGroup *data.BotGroup, bot *data.Bot, stateBase string) (string, int, error) {
	for _, stateLabel := range bot.StateValues {
		// If this is the name of the State we are looking for (BaseName.TargetName)
		if strings.HasPrefix(stateLabel, stateBase+".") {
			index, err := GetStateIndex(botGroup, stateLabel)
			return stateLabel, index, err
		}
	}

	return "", -1, errors.New(fmt.Sprintf("Missing State base: %s  Bot Group: %s  Bot: %s", stateBase, botGroup.Name, bot.Name))
}

// Returns the index of the State inside it's BotForwardSequenceState.  Important because States can only increase or reset to 0 index.
func GetStateIndex(botGroup *data.BotGroup, state string) (int, error) {
	stateSplit := strings.SplitN(state, ".", 2)
	stateBase := stateSplit[0]
	stateLabel := stateSplit[1]

	for _, botState := range botGroup.States {
		if botState.Name == stateBase {
			index, err := util.StringSliceFindIndex(botState.Labels, stateLabel)
			if util.Check(err) {
				return -1, err
			}

			return index, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Missing State: %s  Bot Group: %s", state, botGroup.Name))
}

// Get a BotLockTimer from the BotGroup
func GetLockTimer(botGroup *data.BotGroup, lockTimerName string) (*data.BotLockTimer, error) {
	for _, lockTimer := range botGroup.LockTimers {
		if lockTimer.Name == lockTimerName {
			// If the lock timer is active, but it has timed out, then set to inactive
			if lockTimer.IsActive && lockTimer.Timeout.Unix() < util.GetTimeNow().Unix() {
				lockTimer.IsActive = false
			}

			return &lockTimer, nil
		}
	}

	return &data.BotLockTimer{}, errors.New(fmt.Sprintf("Lock Timer Not Found: %s  Bot Group: %s", lockTimerName, botGroup.Name))
}

func SetLockTimer(botGroup *data.BotGroup, lockTimerName string, duration data.Duration) {
	for _, lockTimer := range botGroup.LockTimers {
		if lockTimer.Name == lockTimerName {
			lockTimer.IsActive = true
			lockTimer.Timeout = util.GetTimeNow().Add(time.Duration(duration))
			return
		}
	}
}

// Get a Bot from the BotGroup
func GetBot(botGroup *data.BotGroup, botName string) (data.Bot, error) {
	for _, bot := range botGroup.Bots {
		if bot.Name == botName {
			return bot, nil
		}
	}
	return data.Bot{}, errors.New(fmt.Sprintf("Bot Group: %s  Bot Missing: %s", botGroup.Name, botName))
}

// Gets a BotGroup from the Site, using the InteractiveControl
func GetBotGroup(interactiveControl data.InteractiveControl, site *data.Site, botGroupName string) (data.BotGroup, error) {
	session := GetInteractiveSession(interactiveControl, site)

	for _, botGroup := range session.BotGroups {
		if botGroup.Name == botGroupName {
			return botGroup, nil
		}
	}

	return data.BotGroup{}, errors.New(fmt.Sprintf("Bot Group Missing: %s", botGroupName))
}

// Returns a slice of Bots in this BotGroup that have this state
func GetBotsInState(botGroup *data.BotGroup, stateName string, stateLabel string) []data.Bot {
	stateKey := fmt.Sprintf("%s.%s", stateName, stateLabel)

	bots := []data.Bot{}

	for _, bot := range botGroup.Bots {
		if util.StringInSlice(bot.StateValues, stateKey) {
			bots = append(bots, bot)
		}
	}

	return bots
}

// Returns all the ConditionCommandResults for all Bots in the BotGroups in this Session, sorted by time, descending
func GetCommandHistoryAll(session *data.InteractiveSession, count int) []data.ConditionCommandResult {
	history := []data.ConditionCommandResult{}

	for _, botGroup := range session.BotGroups {
		for _, bot := range botGroup.Bots {
			for _, commandResult := range bot.CommandHistory {
				history = append(history, commandResult)
			}
		}
	}

	// Sort the list
	sort.Slice(history, fixgo.SliceReverse(func(i, j int) bool { return history[i].Started.Before(history[j].Started) }))

	// Slice from top
	if count > 0 && len(history) > count {
		history = history[0:count]
	}

	return history
}

// ADMIN: Clear the Command History to make the demo look nicer.  Only available if Demo is enabled
func AdminClearCommandHistory() {
	if !data.SireusData.AppConfig.EnableDemo {
		log.Printf("Clearing the command history is not allowed when the Demo is not active.  It is not suitable for a production use case.")
	}

	session := GetInteractiveSession(data.SireusData.Site.ProductionControl, &data.SireusData.Site)
	for botGroupIndex := range session.BotGroups {
		for botIndex := range session.BotGroups[botGroupIndex].Bots {
			bot := &session.BotGroups[botGroupIndex].Bots[botIndex]
			bot.CommandHistory = []data.ConditionCommandResult{}
		}
	}
}

// Returns an array of maps, formatted with the BotVariable name and value
func GetBotGroupAllBotVariablesByName(botGroup data.BotGroup, varName string) []map[string]string {
	varMaps := []map[string]string{}

	for _, bot := range botGroup.Bots {
		for _, value := range bot.SortedVariableValues {
			if value.Key == varName {
				varMap := map[string]string{
					"bot_name":       bot.Name,
					"bot_group_name": botGroup.Name,
					"name":           value.Key,
					"value":          value.Formatted,
					"value_raw":      fmt.Sprintf("%.2f", value.Value),
				}
				varMaps = append(varMaps, varMap)
			}
		}
	}
	return varMaps
}

// Returns a BotGroup variable, by name
func GetBotVariableData(botGroup *data.BotGroup, varName string) (data.BotVariable, error) {
	for _, varData := range botGroup.Variables {
		if varData.Name == varName {
			return varData, nil
		}
	}
	return data.BotVariable{}, errors.New(fmt.Sprintf("Missing variable: %s  Bot Group: %s", varName, botGroup.Name))
}
