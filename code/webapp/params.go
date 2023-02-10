package webapp

import (
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"time"
)

// This is the function that passes in all the data for a given Handlebars page render, using Fiber
func GetPageMapData(c *fiber.Ctx, site *data.Site) fiber.Map {
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

	pageDataMap := BuildRenderMap(site, botGroup, bot)

	return pageDataMap
}

func BuildRenderMap(site *data.Site, botGroup data.BotGroup, bot data.Bot) fiber.Map {
	// Format the Render Time string.  If the Query Time is different, show both so the user knows when they got the
	// information (page load), and when the information query was, if different
	//TODO(ghowland): This will be updated to when we want it to be
	renderTimeStr := util.FormatTimeLong(time.Now())

	renderMap := fiber.Map{
		"title":        "Sireus",
		"site":         site,
		"site_id":      site.Name,
		"botGroup":     botGroup,
		"bot_group_id": botGroup.Name,
		"bot":          bot,
		"bot_id":       bot.Name,
		"render_time":  renderTimeStr,
	}

	return renderMap
}
