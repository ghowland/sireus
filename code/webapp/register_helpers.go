package webapp

import (
	"fmt"
	"github.com/BenJetson/humantime"
	"github.com/aymerick/raymond"
	"github.com/dustin/go-humanize"
	"github.com/ghowland/sireus/code/app"
	"github.com/ghowland/sireus/code/data"
	"github.com/ghowland/sireus/code/util"
	"log"
	"sort"
	"strings"
	"time"
)

// Main function to register all the different Handlebars helper functions, for text processing
func RegisterHandlebarsHelpers() {
	// Testing Length of Arrays for the different structs
	RegisterHandlebarsHelpers_IfArrayLength()

	// Format data
	RegisterHandlebarsHelpers_FormatData()

	// Get AppData values
	RegisterHandlebarsHelpers_GetAppData()

	// Get Go data values
	RegisterHandlebarsHelpers_GetGoData()

	// Sets current data from otherwise inaccessible data structures, because of slicing, map references, looks ups, etc
	RegisterHandlebarsHelpers_WithData()

	// Expanded test logic
	RegisterHandlebarsHelpers_IfTests()
}

// Sets current data from otherwise inaccessible data structures, because of slicing, map references, looks ups, etc
func RegisterHandlebarsHelpers_WithData() {
	// With BotGroup by name
	raymond.RegisterHelper("with_bot_group_by_name", func(session data.InteractiveSession, name string, options *raymond.Options) raymond.SafeString {
		for _, botGroup := range session.BotGroups {
			if botGroup.Name == name {
				return raymond.SafeString(options.FnWith(botGroup))
			}
		}
		return raymond.SafeString(options.FnWith(data.BotGroup{Name: "Missing"}))
	})

	// With Bots in specified state
	raymond.RegisterHelper("with_bots_in_state", func(botGroup data.BotGroup, stateName string, stateLabel string, options *raymond.Options) raymond.SafeString {
		bots := app.GetBotsInState(&botGroup, stateName, stateLabel)
		return raymond.SafeString(options.FnWith(bots))
	})

	// With Count of Bots in specified state
	raymond.RegisterHelper("with_command_history_all_latest", func(session data.InteractiveSession, count int, options *raymond.Options) raymond.SafeString {
		allCommandHistory := app.GetCommandHistoryAll(&session, count)
		return raymond.SafeString(options.FnWith(allCommandHistory))
	})

	// With BotConditionData
	raymond.RegisterHelper("with_bot_condition", func(bot data.Bot, condition data.Condition, options *raymond.Options) raymond.SafeString {
		botConditionData := bot.ConditionData[condition.Name]
		return raymond.SafeString(options.FnWith(botConditionData))
	})

	// With Condition from Bot
	raymond.RegisterHelper("with_condition_from_bot", func(botGroup data.BotGroup, bot data.Bot, conditionDataIndex int, options *raymond.Options) raymond.SafeString {

		botConditionData := bot.SortedConditionData[conditionDataIndex]

		botCondition, err := app.GetCondition(&botGroup, botConditionData.Key)
		util.CheckLog(err)
		return raymond.SafeString(options.FnWith(botCondition))
	})

	// With Query Server by Name from Site
	raymond.RegisterHelper("with_query_server", func(queryServerName string, site data.Site, options *raymond.Options) raymond.SafeString {
		queryServer, err := app.GetQueryServer(&site, queryServerName)
		util.CheckLog(err)

		return raymond.SafeString(options.FnWith(queryServer))
	})

	// With Query Server by Name from Site
	raymond.RegisterHelper("with_bot_group_bot_variable_by_name", func(botGroup data.BotGroup, varName string, options *raymond.Options) raymond.SafeString {
		variables := app.GetBotGroupAllBotVariablesByName(botGroup, varName)
		return raymond.SafeString(options.FnWith(variables))
	})
}

// Expanded test logic
func RegisterHandlebarsHelpers_IfTests() {
	// If string == string
	raymond.RegisterHelper("if_equal_string", func(a string, b string, options *raymond.Options) raymond.SafeString {
		if a == b {
			log.Printf("Equal String: True: %s == %s -> %v", a, b, a == b)
			return raymond.SafeString(options.Fn())
		} else {
			log.Printf("Equal String: False: %s == %s -> %v", a, b, a == b)
			return raymond.SafeString("")
		}
	})

	// If string in []string
	raymond.RegisterHelper("if_string_in_slice", func(slice []string, find string, options *raymond.Options) interface{} {
		if util.StringInSlice(slice, find) {
			return raymond.SafeString(options.Fn())
		} else {
			return options.Inverse()
		}
	})

	// If []string has concatenation of strings in it
	raymond.RegisterHelper("if_slice_has_dot_strings_2", func(a []string, b1 string, b2 string, options *raymond.Options) interface{} {
		testString := fmt.Sprintf("%s.%s", b1, b2)

		if util.StringInSlice(a, testString) {
			return raymond.SafeString(options.Fn())
		} else {
			return options.Inverse()
		}
	})

	// If Go time.Time == 0
	raymond.RegisterHelper("if_time_never", func(t time.Time, options *raymond.Options) raymond.SafeString {
		if t.IsZero() {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	// If Go time.Time != 0
	raymond.RegisterHelper("if_not_time_never", func(t time.Time, options *raymond.Options) raymond.SafeString {
		if !t.IsZero() {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})
}

// Get Go data values.  Slices, maps, etc
func RegisterHandlebarsHelpers_GetGoData() {
	// Consideration Scores: Final
	raymond.RegisterHelper("get_string_slice_index", func(stringArray []string, index int) raymond.SafeString {
		value := fmt.Sprintf("MISSING:%d", index)
		// If it's negative, fix it to be positive
		if index < 0 {
			index = len(stringArray) + index
		}
		// Test if it's a valid index, otherwise we return the missing pre-set string
		if index > 0 && index < len(stringArray) {
			value = stringArray[index]
		}
		return raymond.SafeString(value)
	})
}

// Get AppData values.  Bot, BotGroup, Condition, BotConditionData, etc
func RegisterHandlebarsHelpers_GetAppData() {
	// Consideration Scores: Final
	raymond.RegisterHelper("get_bot_condition_data_consideration_final_score", func(bot data.Bot, condition data.Condition, consider data.ConditionConsideration) raymond.SafeString {
		util.LockAcquire(bot.LockKey)
		defer util.LockRelease(bot.LockKey)
		output := fmt.Sprintf("%.2f", bot.ConditionData[condition.Name].ConsiderationFinalScores[consider.Name])
		return raymond.SafeString(output)
	})

	// Consideration Scores: Raw (not Ranged, Curved, Weighted)
	raymond.RegisterHelper("get_bot_condition_data_consideration_raw_score", func(bot data.Bot, condition data.Condition, consider data.ConditionConsideration) raymond.SafeString {
		util.LockAcquire(bot.LockKey)
		defer util.LockRelease(bot.LockKey)
		output := fmt.Sprintf("%.2f", bot.ConditionData[condition.Name].ConsiderationRawScores[consider.Name])
		return raymond.SafeString(output)
	})

	// Consideration Scores: Ranged (not Curved, Weighted)
	raymond.RegisterHelper("get_bot_condition_data_consideration_ranged_score", func(bot data.Bot, condition data.Condition, consider data.ConditionConsideration) raymond.SafeString {
		util.LockAcquire(bot.LockKey)
		defer util.LockRelease(bot.LockKey)
		output := fmt.Sprintf("%.2f", bot.ConditionData[condition.Name].ConsiderationRangedScores[consider.Name])
		return raymond.SafeString(output)
	})

	// Consideration Scores: Curved (not Weighted)
	raymond.RegisterHelper("get_bot_condition_data_consideration_curved_score", func(bot data.Bot, condition data.Condition, consider data.ConditionConsideration) raymond.SafeString {
		util.LockAcquire(bot.LockKey)
		defer util.LockRelease(bot.LockKey)
		output := fmt.Sprintf("%.2f", bot.ConditionData[condition.Name].ConsiderationCurvedScores[consider.Name])
		return raymond.SafeString(output)
	})

	// ConditionData Final Score
	raymond.RegisterHelper("get_bot_condition_data_final_score", func(bot data.Bot, condition data.Condition) raymond.SafeString {
		util.LockAcquire(bot.LockKey)
		defer util.LockRelease(bot.LockKey)
		output := fmt.Sprintf("%.2f", bot.ConditionData[condition.Name].FinalScore)
		return raymond.SafeString(output)
	})

	// Get Count of Bot slice
	raymond.RegisterHelper("get_len_array_bot", func(bots []data.Bot, options *raymond.Options) raymond.SafeString {
		return raymond.SafeString(options.FnWith(len(bots)))
	})

	// Get a slice of names from a slice of Bots
	raymond.RegisterHelper("get_array_bot_names", func(bots []data.Bot, options *raymond.Options) raymond.SafeString {
		botNames := []string{}
		for _, bot := range bots {
			botNames = append(botNames, bot.Name)
		}
		sort.Strings(botNames)
		return raymond.SafeString(options.FnWith(botNames))
	})
}

// Format data, for Go and our internal data types
func RegisterHandlebarsHelpers_FormatData() {
	// Queries
	raymond.RegisterHelper("format_query_web", func(site data.Site, item data.BotQuery) string {
		queryServer, err := app.GetQueryServer(&site, item.QueryServer)
		util.CheckLog(err)
		mapData := map[string]string{
			"query": item.Query,
		}
		return util.HandlebarFormatText(queryServer.WebUrlFormat, mapData)
	})

	// Queries
	raymond.RegisterHelper("format_query_server_web", func(queryServer data.QueryServer, query string) string {
		mapData := map[string]string{
			"query": query,
		}
		return util.HandlebarFormatText(queryServer.WebUrlFormat, mapData)
	})

	// Format Go Values
	raymond.RegisterHelper("format_float64", func(format string, value float64) raymond.SafeString {
		output := fmt.Sprintf(format, value)
		//log.Printf("Format: %s  Val: %v  Output: %v", format, value, output)
		return raymond.SafeString(output)
	})

	// Format Time
	raymond.RegisterHelper("format_time_since", func(t time.Time) raymond.SafeString {
		return raymond.SafeString(humantime.Since(t))
	})

	raymond.RegisterHelper("format_time_since_precise", func(t time.Time) raymond.SafeString {
		return raymond.SafeString(humanize.Time(t))
	})

	raymond.RegisterHelper("format_time", func(t time.Time) raymond.SafeString {
		return raymond.SafeString(util.FormatTimeLong(t))
	})

	raymond.RegisterHelper("format_duration", func(d data.Duration) raymond.SafeString {
		return raymond.SafeString(time.Duration(d).String())
	})

	raymond.RegisterHelper("format_html_id", func(name string) raymond.SafeString {
		nameId := util.StringReplaceUnsafeChars(name, " [](){}=-!@#$%^&*()+<>,./?;:'\"`~", "_")
		return raymond.SafeString(nameId)
	})

	// Format the Formatting
	raymond.RegisterHelper("format_variable_format", func(item data.BotVariableFormat) string {
		return item.String()
	})

	// []string
	raymond.RegisterHelper("format_array_string_csv", func(item []string) string {
		return strings.Join(item, ", ")
	})

	// Format a substring slice of this string
	raymond.RegisterHelper("format_string_substring", func(value string, start int, end int) string {
		if len(value) == 0 {
			return ""
		}

		if end >= len(value) {
			end = len(value)
		}
		if start >= end {
			start = end - 1
		}
		if start < 0 {
			start = 0
		}
		return value[start:end]
	})

	// Formats a config file, so we can show it to the user.  Comes our from HTML Handlebars
	raymond.RegisterHelper("format_config_file", func(path string) string {
		// Don't allow pathing operations
		if strings.HasPrefix(path, ".") {
			return fmt.Sprintf("Forbidden: %s", path)
		} else if strings.HasPrefix(path, "/") {
			return fmt.Sprintf("Forbidden: %s", path)
		} else if strings.HasPrefix(path, "~") {
			return fmt.Sprintf("Forbidden: %s", path)
		}

		// Don't allow backtracking or any modification of path, or looking for "hidden files"
		if strings.Contains(path, "/.") {
			return fmt.Sprintf("Forbidden: %s", path)
		}

		content, err := util.FileLoad(path)
		if util.Check(err) {
			return fmt.Sprintf("Missing path: %s", path)
		}
		return content
	})

}

// Testing Length of Arrays for the different structs
func RegisterHandlebarsHelpers_IfArrayLength() {
	// -- Go Data --

	// The data structure needs to be []interface{} to work, it wont auto-cast from Handlerbars to here, like []app.Bots -> []interface{}
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
	raymond.RegisterHelper("if_bot_group_length", func(items []data.BotGroup, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_bot_length", func(items []data.Bot, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_condition_length", func(items []data.Condition, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_command_history_length", func(items []data.ConditionCommandResult, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_consider_length", func(items []data.ConditionConsideration, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_query_length", func(items []data.BotQuery, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_variable_length", func(items []data.BotVariable, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_state_length", func(items []data.BotForwardSequenceState, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})

	raymond.RegisterHelper("if_lock_timer_length", func(items []data.BotLockTimer, count int, options *raymond.Options) raymond.SafeString {
		if len(items) >= count {
			return raymond.SafeString(options.Fn())
		} else {
			return raymond.SafeString("")
		}
	})
}
