package webapp

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/ghowland/sireus/code/appdata"
	"github.com/ghowland/sireus/code/util"
	"strings"
)

func RegisterHandlebarsHelpers() {
	// Testing Length of Arrays for the different structs
	RegisterHandlebarsHelpers_IfArrayLength()

	// Format data
	RegisterHandlebarsHelpers_FormatData()

	// Get AppData values
	RegisterHandlebarsHelpers_GetAppData()

	// Get values from structs.  Only needed for when data isn't available as a child.
	RegisterHandlebarsHelpers_GetValuesFromStructs()
}

func RegisterHandlebarsHelpers_GetValuesFromStructs() {
	// BotActionData
	raymond.RegisterHelper("struct_bot_action_data", func(bot appdata.Bot, action appdata.Action) appdata.BotActionData {
		botActionData := bot.ActionData[action.Name]
		return botActionData
	})

}

func RegisterHandlebarsHelpers_GetAppData() {
	// Consideration Scores: Final
	raymond.RegisterHelper("get_bot_action_data_consideration_final_score", func(bot appdata.Bot, action appdata.Action, consider appdata.ActionConsideration) raymond.SafeString {
		return raymond.SafeString(fmt.Sprintf("%.2f", bot.ActionData[action.Name].ConsiderationFinalScores[consider.Name]))
	})

	// Consideration Scores: Calculated (not Weighted)
	raymond.RegisterHelper("get_bot_action_data_consideration_evaluated_score", func(bot appdata.Bot, action appdata.Action, consider appdata.ActionConsideration) raymond.SafeString {
		return raymond.SafeString(fmt.Sprintf("%.2f", bot.ActionData[action.Name].ConsiderationEvaluatedScores[consider.Name]))
	})

	// ActionData Final Score
	raymond.RegisterHelper("get_bot_action_data_final_score", func(bot appdata.Bot, action appdata.Action) raymond.SafeString {
		return raymond.SafeString(fmt.Sprintf("%.2f", bot.ActionData[action.Name].FinalScore))
	})
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
	raymond.RegisterHelper("format_variable_type", func(item appdata.BotVariableType) string {
		return item.String()
	})

	// []string
	raymond.RegisterHelper("format_array_string_csv", func(item []string) string {
		return strings.Join(item, ", ")
	})
}

func RegisterHandlebarsHelpers_IfArrayLength() {
	// -- Go Data --

	// The data structure needs to be []interface{} to work, it wont auto-cast from Handlerbars to here, like []appdata.Bots -> []interface{}
	raymond.RegisterHelper("if_array_length", func(items []interface{}, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_map_string_float64_length", func(items map[string]float64, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	// -- Sireus Structs --

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
