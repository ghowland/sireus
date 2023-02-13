package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"os"
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
		site.LoadedBotGroups = append(site.LoadedBotGroups, botGroup)
	}

	return site
}

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
	session.TimeRequested = time.Now()
	session.QueryStartTime = time.UnixMilli(int64(interactiveControl.QueryStartTime))
	session.QueryDuration = data.Duration(interactiveControl.QueryDuration)
	session.QueryScrubTime = time.UnixMilli(int64(interactiveControl.QueryScrubTime))

	// If this is a Production session, we invalidate cache on interval and dont worry about query mismatch
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
		QueryStartTime:         float64(time.Now().UnixMilli()),
		QueryDuration:          60 * 1000000000,
		QueryScrubTime:         float64(time.Now().UnixMilli()),
	}

	return interactiveControl
}

// Returns a QueryServer, scope is per Site
func GetQueryServer(site data.Site, name string) (data.QueryServer, error) {
	for _, queryServer := range site.QueryServers {
		if queryServer.Name == name {
			return queryServer, nil
		}
	}

	return data.QueryServer{}, errors.New(fmt.Sprintf("Query Server missing: %s", name))
}

// Gets a query, scope per BotGroup
func GetQuery(botGroup data.BotGroup, queryName string) (data.BotQuery, error) {
	for _, query := range botGroup.Queries {
		if query.Name == queryName {
			return query, nil
		}
	}
	return data.BotQuery{}, errors.New(fmt.Sprintf("Bot Group: %s  Query missing: %s", botGroup.Name, queryName))
}

// Get an Action from a BotGroup, by name
func GetAction(botGroup data.BotGroup, actionName string) (data.Action, error) {
	for _, action := range botGroup.Actions {
		if action.Name == actionName {
			return action, nil
		}
	}
	return data.Action{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Action: %s", botGroup.Name, actionName))
}

// Get a Variable defintion from BotGroup, by name.  Not the Variable Value, which is stored in Bot.
func GetVariable(botGroup data.BotGroup, varName string) (data.BotVariable, error) {
	for _, variable := range botGroup.Variables {
		if variable.Name == varName {
			return variable, nil
		}
	}
	return data.BotVariable{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Variable: %s", botGroup.Name, varName))
}

// Get an ActionConsideration from an Action, by name
func GetActionConsideration(action data.Action, considerName string) (data.ActionConsideration, error) {
	for _, consider := range action.Considerations {
		if consider.Name == considerName {
			return consider, nil
		}
	}
	return data.ActionConsideration{}, errors.New(fmt.Sprintf("Missing Consideration: %s", considerName))
}

// For a given Action, does this Bot have all the RequiredStates active?
func AreAllActionStatesActive(action data.Action, bot data.Bot) bool {
	for _, state := range action.RequiredStates {
		if !util.StringInSlice(state, bot.StateValues) {
			return false
		}
	}

	return true
}

// Get a Bot from the BotGroup
func GetBot(botGroup data.BotGroup, botName string) (data.Bot, error) {
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
