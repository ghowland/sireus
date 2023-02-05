package webapp

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
)

func RegisterHandlebarsHelpers() {
	//TEST:REMOVE: Not actually useful, just demo.  Remove once I have a real version of this that is a better demo
	//NOTE(ghowalnd): This demo could use Partials for rendering too, so it's not hard coding anything
	raymond.RegisterHelper("botinfo", func(bot appdata.BotGroup) string {
		return fmt.Sprintf("%s  Actions: %d", bot.Name, len(bot.Actions))
	})

	// Testing Length of Arrays for the different structs
	RegisterHandlebarsHelpers_IfArrayLength()

	// Format data
	RegisterHandlebarsHelpers_FormatData()
}

func RegisterHandlebarsHelpers_FormatData() {
	// Queries
	raymond.RegisterHelper("format_query_web", func(site appdata.Site, item appdata.BotQuery) string {
		queryServer, err := appdata.GetQueryServer(site, item.QueryServer)
		util.Check(err)
		mapData := map[string]string{
			"query": item.Query,
		}
		return util.HandlebarFormatText(queryServer.WebUrlFormat, mapData)
	})

	// Variables
	raymond.RegisterHelper("format_variable_type", func(item appdata.BotVariable) string {
		return item.Type.String()
	})
}

func RegisterHandlebarsHelpers_IfArrayLength() {
	// NOTE(ghowland): I am choosing to do this on a per-data type basis instead of generalizing, as it will make
	//				   targeted changes faster and easier in the future

	// Testing Length of Arrays for the different structs
	raymond.RegisterHelper("if_bot_group_length", func(items []appdata.BotGroup, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_bot_length", func(items []appdata.Bot, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_action_length", func(items []appdata.Action, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_consider_length", func(items []appdata.ActionConsideration, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_query_length", func(items []appdata.BotQuery, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_variable_length", func(items []appdata.BotVariable, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_state_length", func(items []appdata.BotForwardSequenceState, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_lock_timer_length", func(items []appdata.BotLockTimer, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})
}
