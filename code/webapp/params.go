package webapp

import (
	"encoding/json"
	"fmt"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
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

// GetRenderMapFromRPC parses RPC params and passes in all the data for a given Handlebars page render, using go map.
// This uses the Interactive Control data to modify data based on the settings.
func GetRenderMapFromRPC(c *fiber.Ctx, siteActual *data.Site) map[string]interface{} {
	input := util.ParseContextBody(c)

	// Get our Interactive Controls, if it exists
	var interactiveControl data.InteractiveControl
	interactiveControlJSON, ok := input["interactive_control"]
	if ok {
		json.Unmarshal([]byte(interactiveControlJSON), &interactiveControl)
	}
	log.Printf("RPC Args: Interactive: %s", util.PrintJson(interactiveControl))

	if interactiveControl.SessionUUID == 0 {
		interactiveControl.SessionUUID = uuid.New().ID()
	}

	botGroupId := input["bot_group_id"]
	botId := input["bot_id"]

	// If we are using customized data, then we make a new site variable
	site := siteActual
	if interactiveControl.UseInteractiveSession {
		site = siteActual //TODO: Change this to a new version of site, cloned off of this, but with our Interactive data
	}

	// Bot Groups and Bots come from the Site.  Site is either original or the Interactive data version, but treated the same
	botGroup := data.BotGroup{}
	var err error
	if len(botGroupId) != 0 {
		botGroup, err = app.GetBotGroup(site, botGroupId)
		util.Check(err)
	}

	bot := data.Bot{}
	if len(botId) != 0 && len(botGroup.Name) != 0 {
		bot, err = app.GetBot(botGroup, botId)
		util.Check(err)
	}

	inputData := make(map[string]interface{})
	inputData["bot_group_id"] = botGroupId
	inputData["bot_id"] = botId

	// The site will remain the same, because it also has all our queries and lock timers and everything else.
	renderMap := BuildRenderMap(site, botGroup, bot, inputData, interactiveControl)

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

	interactiveStartTime := FormatInteractiveStartTime()

	renderMap := fiber.Map{
		"title":                    "Sireus",
		"site":                     site,
		"site_id":                  site.Name,
		"botGroup":                 botGroup,
		"bot_group_id":             botGroup.Name,
		"bot":                      bot,
		"bot_id":                   bot.Name,
		"render_time":              renderTimeStr,
		"input_data":               inputDataStr,
		"interactive_starter_time": interactiveStartTime,
		"interactive_control":      "{}", // Always empty from initial page render
	}

	return renderMap
}

func FormatInteractiveStartTime() string {
	// 15 minutes ago
	//TODO(ghowland): Remove hard-code, put into AppConfig, also make default Duration in the webapp
	var t = time.Now().Add(-15 * 60 * time.Second).UTC()

	ampm := "AM"
	hour := t.Hour()
	if hour > 12 {
		hour -= 12
		ampm = "PM"
	}

	output := fmt.Sprintf("%02d/%02d/%d, %d:%02d %s", t.Day(), t.Month(), t.Year(), hour, t.Minute(), ampm)
	return output
}

func BuildRenderMap(site *data.Site, botGroup data.BotGroup, bot data.Bot, inputData map[string]interface{}, interactiveControl data.InteractiveControl) map[string]interface{} {
	// Format the Render Time string.  If the Query Time is different, show both so the user knows when they got the
	// information (page load), and when the information query was, if different
	//TODO(ghowland): This will be updated to when we want it to be
	renderTimeStr := util.FormatTimeLong(time.Now())

	queryTime := time.UnixMilli(int64(interactiveControl.QueryScrubTime))
	queryTimeStr := util.FormatTimeLong(queryTime)

	inputDataStr := strings.Replace(util.PrintJsonData(inputData), "\"", "\\\"", -1)

	interactiveControlStr := strings.Replace(util.PrintJsonData(interactiveControl), "\"", "\\\"", -1)

	// If we got nothing, pass in empty values, so the JSON is still valid
	if len(inputDataStr) == 0 {
		inputDataStr = "{}"
	}
	if len(interactiveControlStr) == 0 {
		interactiveControlStr = "{}"
	}

	renderMap := map[string]interface{}{
		"title":               "Sireus",
		"site":                site,
		"site_id":             site.Name,
		"botGroup":            botGroup,
		"bot_group_id":        botGroup.Name,
		"bot":                 bot,
		"bot_id":              bot.Name,
		"render_time":         renderTimeStr,
		"query_time":          queryTimeStr,
		"input_data":          inputDataStr,
		"interactive_control": interactiveControlStr,
	}

	return renderMap
}
