package webapp

import (
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// GetRenderMapFromParams parses GET params passes in all the data for a given Handlebars page render, using Fiber
func GetRenderMapFromParams(c *fiber.Ctx, site *data.Site) fiber.Map {
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

	inputData := make(map[string]interface{})
	inputData["bot_group_id"] = botGroupId
	inputData["bot_id"] = botId

	renderMap := BuildRenderMapFiber(site, botGroup, bot, inputData)

	return renderMap
}

// GetRenderMapFromRPC parses RPC params and passes in all the data for a given Handlebars page render, using go map
func GetRenderMapFromRPC(c *fiber.Ctx, site *data.Site) map[string]interface{} {
	input := util.ParseContextBody(c)

	botGroupId := input["bot_group_id"]
	botId := input["bot_id"]

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

	inputData := make(map[string]interface{})
	inputData["bot_group_id"] = botGroupId
	inputData["bot_id"] = botId

	renderMap := BuildRenderMap(site, botGroup, bot, inputData)

	return renderMap
}

func BuildRenderMapFiber(site *data.Site, botGroup data.BotGroup, bot data.Bot, inputData map[string]interface{}) fiber.Map {
	// Format the Render Time string.  If the Query Time is different, show both so the user knows when they got the
	// information (page load), and when the information query was, if different
	//TODO(ghowland): This will be updated to when we want it to be
	renderTimeStr := util.FormatTimeLong(time.Now())

	inputDataStr := strings.Replace(util.PrintJsonData(inputData), "\"", "\\\"", -1)

	// If we got nothing, pass in empty values, so the JSON is still valid
	if len(inputData) == 0 {
		inputDataStr = "{}"
	}

	renderMap := fiber.Map{
		"title":        "Sireus",
		"site":         site,
		"site_id":      site.Name,
		"botGroup":     botGroup,
		"bot_group_id": botGroup.Name,
		"bot":          bot,
		"bot_id":       bot.Name,
		"render_time":  renderTimeStr,
		"input_data":   inputDataStr,
	}

	return renderMap
}

func BuildRenderMap(site *data.Site, botGroup data.BotGroup, bot data.Bot, inputData map[string]interface{}) map[string]interface{} {
	// Format the Render Time string.  If the Query Time is different, show both so the user knows when they got the
	// information (page load), and when the information query was, if different
	//TODO(ghowland): This will be updated to when we want it to be
	renderTimeStr := util.FormatTimeLong(time.Now())

	inputDataStr := strings.Replace(util.PrintJsonData(inputData), "\"", "\\\"", -1)

	// If we got nothing, pass in empty values, so the JSON is still valid
	if len(inputData) == 0 {
		inputDataStr = "{}"
	}

	renderMap := map[string]interface{}{
		"title":        "Sireus",
		"site":         site,
		"site_id":      site.Name,
		"botGroup":     botGroup,
		"bot_group_id": botGroup.Name,
		"bot":          bot,
		"bot_id":       bot.Name,
		"render_time":  renderTimeStr,
		"input_data":   inputDataStr,
	}

	return renderMap
}
