package webapp

import (
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
)

func GetPageMapData(c *fiber.Ctx, site appdata.Site) fiber.Map {
	siteId := site.Name
	botGroupId := c.Query("bot_group_id")
	botId := c.Query("bot_id")

	botGroup := appdata.BotGroup{}
	var err error
	if botGroupId != "" {
		botGroup, err = appdata.GetBotGroup(site, botGroupId)
		util.Check(err)
	}

	bot := appdata.Bot{}
	if botId != "" && botGroup.Name != "" {
		bot, err = appdata.GetBot(botGroup, botId)
		util.Check(err)
	}

	pageDataMap := fiber.Map{
		"title":        "Sireus",
		"site":         site,
		"botGroup":     botGroup,
		"site_id":      siteId,
		"bot_group_id": botGroupId,
		"bot_id":       botId,
		"bot":          bot,
	}

	return pageDataMap
}
