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
	util.Check(err)

	var botGroup BotGroup
	err = json.Unmarshal(botGroupData, &botGroup)
	util.Check(err)

	return botGroup
}

func LoadSiteConfig(appConfig AppConfig) Site {
	siteData, err := os.ReadFile(appConfig.SiteConfigPath)
	util.Check(err)

	var site Site
	err = json.Unmarshal(siteData, &site)
	util.Check(err)

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
