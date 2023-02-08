package webapp

import (
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
)

// This is the function that passes in all the data for a given Handlebars page render, using Fiber
func GetPageMapData(c *fiber.Ctx, site data.Site) fiber.Map {
	siteId := site.Name
	botGroupId := c.Query("bot_group_id")
	botId := c.Query("bot_id")

	botGroup := data.BotGroup{}
	var err error
	if botGroupId != "" {
		botGroup, err = app.GetBotGroup(site, botGroupId)
		util.Check(err)
	}

	bot := data.Bot{}
	if botId != "" && botGroup.Name != "" {
		bot, err = app.GetBot(botGroup, botId)
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
