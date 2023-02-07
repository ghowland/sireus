package appdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghowland/sireus/code/util"
	"os"
)

func LoadBotGroupConfig(path string) BotGroup {
	botGroupData, err := os.ReadFile(path)
	util.CheckPanic(err)

	var botGroup BotGroup
	err = json.Unmarshal(botGroupData, &botGroup)
	util.CheckPanic(err)

	return botGroup
}

func LoadSiteConfig(appConfig AppConfig) Site {
	siteData, err := os.ReadFile(appConfig.SiteConfigPath)
	util.CheckPanic(err)

	var site Site
	err = json.Unmarshal(siteData, &site)
	util.CheckPanic(err)

	// Load all our Bot Groups
	for _, botGroupPath := range site.BotGroupPaths {
		botGroup := LoadBotGroupConfig(botGroupPath)
		site.BotGroups = append(site.BotGroups, botGroup)
	}

	return site
}

func GetQueryServer(site Site, name string) (QueryServer, error) {
	for _, queryServer := range site.QueryServers {
		if queryServer.Name == name {
			return queryServer, nil
		}
	}

	return QueryServer{}, errors.New(fmt.Sprintf("Query Server missing: %s", name))
}

func GetQuery(botGroup BotGroup, queryName string) (BotQuery, error) {
	for _, query := range botGroup.Queries {
		if query.Name == queryName {
			return query, nil
		}
	}
	return BotQuery{}, errors.New(fmt.Sprintf("Bot Group: %s  Query missing: %s", botGroup.Name, queryName))
}

func GetBotGroup(site Site, botGroupName string) (BotGroup, error) {
	for _, botGroup := range site.BotGroups {
		if botGroup.Name == botGroupName {
			return botGroup, nil
		}
	}
	return BotGroup{}, errors.New(fmt.Sprintf("Bot Ground Missing: %s", botGroupName))
}

func GetBot(site Site, botGroup BotGroup, botName string) (Bot, error) {
	for _, bot := range botGroup.Bots {
		if bot.Name == botName {
			return bot, nil
		}
	}
	return Bot{}, errors.New(fmt.Sprintf("Bot Group: %s  Bot Missing: %s", botGroup.Name, botName))
}

func GetAction(botGroup BotGroup, actionName string) (Action, error) {
	for _, action := range botGroup.Actions {
		if action.Name == actionName {
			return action, nil
		}
	}
	return Action{}, errors.New(fmt.Sprintf("Bot Group: %s  Missing Action: %s", botGroup.Name, actionName))
}

func GetActionConsideration(action Action, considerName string) (ActionConsideration, error) {
	for _, consider := range action.Considerations {
		if consider.Name == considerName {
			return consider, nil
		}
	}
	return ActionConsideration{}, errors.New(fmt.Sprintf("Missing Consideration: %s", considerName))
}
