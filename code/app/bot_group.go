package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"os"
)

func LoadBotGroupConfig(path string) data.BotGroup {
	botGroupData, err := os.ReadFile(path)
	util.CheckPanic(err)

	var botGroup data.BotGroup
	err = json.Unmarshal(botGroupData, &botGroup)
	util.CheckPanic(err)

	return botGroup
}

func LoadSiteConfig(appConfig data.AppConfig) data.Site {
	siteData, err := os.ReadFile(appConfig.SiteConfigPath)
	util.CheckPanic(err)

	var site data.Site
	err = json.Unmarshal(siteData, &site)
	util.CheckPanic(err)

	// Load all our Bot Groups
	for _, botGroupPath := range site.BotGroupPaths {
		botGroup := LoadBotGroupConfig(botGroupPath)
		site.BotGroups = append(site.BotGroups, botGroup)
	}

	return site
}

func GetQueryServer(site data.Site, name string) (data.QueryServer, error) {
	for _, queryServer := range site.QueryServers {
		if queryServer.Name == name {
			return queryServer, nil
		}
	}

	return data.QueryServer{}, errors.New(fmt.Sprintf("Query Server missing: %s", name))
}

func GetQuery(botGroup data.BotGroup, queryName string) (data.BotQuery, error) {
	for _, query := range botGroup.Queries {
		if query.Name == queryName {
			return query, nil
		}
	}
	return data.BotQuery{}, errors.New(fmt.Sprintf("Bot Group: %s  Query missing: %s", botGroup.Name, queryName))
}

func GetBotGroup(site data.Site, botGroupName string) (data.BotGroup, error) {
	for _, botGroup := range site.BotGroups {
		if botGroup.Name == botGroupName {
			return botGroup, nil
		}
	}
	return data.BotGroup{}, errors.New(fmt.Sprintf("Bot Ground Missing: %s", botGroupName))
}

func GetBot(botGroup data.BotGroup, botName string) (data.Bot, error) {
	for _, bot := range botGroup.Bots {
		if bot.Name == botName {
			return bot, nil
		}
	}
	return data.Bot{}, errors.New(fmt.Sprintf("Bot Group: %s  Bot Missing: %s", botGroup.Name, botName))
}

func GetAction(botGroup data.BotGroup, actionName string) (data.Action, error) {
	for _, action := range botGroup.Actions {
		if action.Name == actionName {
			return action, nil
		}
	}
	return data.Action{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Action: %s", botGroup.Name, actionName))
}

func GetVariable(botGroup data.BotGroup, varName string) (data.BotVariable, error) {
	for _, variable := range botGroup.Variables {
		if variable.Name == varName {
			return variable, nil
		}
	}
	return data.BotVariable{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Variable: %s", botGroup.Name, varName))
}

func GetActionConsideration(action data.Action, considerName string) (data.ActionConsideration, error) {
	for _, consider := range action.Considerations {
		if consider.Name == considerName {
			return consider, nil
		}
	}
	return data.ActionConsideration{}, errors.New(fmt.Sprintf("Missing Consideration: %s", considerName))
}

func AreAllActionStatesActive(action data.Action, bot data.Bot) bool {
	for _, state := range action.RequiredStates {
		if !util.StringInSlice(state, bot.StateValues) {
			return false
		}
	}

	return true
}